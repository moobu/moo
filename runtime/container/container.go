package container

import (
	"context"
	"sync"
	"time"

	"github.com/moobu/moo/runtime"
)

type Client interface {
	Create(string, ...CreateOption) (*Container, error)
	Inspect(string, ...InspectOption) (*Container, error)
	Delete(string, ...DeleteOption) error
	List(...ListOption) ([]*Container, error)
	Wait(*Container, ...WaitOption) error
	Start(*Container, ...StartOption) error
	Stop(*Container, ...StopOption) error
}

type Container struct {
	ID     string
	Status runtime.Status
}

type container struct {
	sync.RWMutex
	pods    map[string]map[string]*cpod
	options runtime.Options
	client  Client
	wg      sync.WaitGroup
	exit    chan struct{}
}

func (c *container) Create(pod *runtime.Pod, opts ...runtime.CreateOption) error {
	c.Lock()
	defer c.Unlock()

	var options runtime.CreateOptions
	for _, o := range opts {
		o(&options)
	}

	id := pod.String()
	ns := options.Namespace

	if _, ok := c.pods[ns]; !ok {
		c.pods[ns] = make(map[string]*cpod)
	}
	if _, ok := c.pods[id]; ok {
		return runtime.ErrExists
	}

	ctx := options.Context
	con, err := c.client.Create(options.Bundle.Image, CreateContext(ctx))
	if err != nil {
		return err
	}

	p := &cpod{
		Pod:       pod,
		wg:        &c.wg,
		client:    c.client,
		container: con,
		retries:   options.Retries,
	}
	if err := p.start(ctx); err != nil {
		return err
	}

	c.pods[ns][id] = p
	return nil
}

func (c *container) Delete(pod *runtime.Pod, opts ...runtime.DeleteOption) error {
	c.Lock()
	defer c.Unlock()

	var options runtime.DeleteOptions
	for _, o := range opts {
		o(&options)
	}

	id := pod.String()
	ns := options.Namespace
	if _, ok := c.pods[ns]; !ok {
		return nil
	}

	p, ok := c.pods[ns][id]
	if !ok {
		return runtime.ErrNotFound
	}

	ctx := options.Context
	con, err := c.client.Inspect(p.container.ID, InspectContext(ctx))
	if err != nil {
		return err
	}

	if con.Status == runtime.Exited {
		return nil
	}

	if err := p.stop(ctx); err != nil {
		return err
	}
	delete(c.pods[ns], id)
	return nil
}

func (c *container) List(opts ...runtime.ListOption) ([]*runtime.Pod, error) {
	c.Lock()
	defer c.Unlock()

	var options runtime.ListOptions
	for _, o := range opts {
		o(&options)
	}

	ns := options.Namespace
	ps, ok := c.pods[ns]
	if !ok {
		return nil, nil
	}

	// TODO: fetch the container runtime in case of a verbose list
	// ctx := options.Context
	// if options.Verbose {
	// 	label := fmt.Sprintf("namespace=%s", ns)
	// 	cons, err := c.client.List(Labels(label), ListContext(ctx))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	pods := make([]*runtime.Pod, len(ps))
	for _, p := range ps {
		pods = append(pods, p.Pod)
	}
	return pods, nil
}

func (c *container) Start() (err error) {
	c.Lock()
	defer c.Unlock()

	var events <-chan runtime.Event
	if c.options.Scheduler != nil {
		if events, err = c.options.Scheduler.Schedule(); err != nil {
			return
		}
	}

	go c.run(events) // start the runtime daemon
	return
}

func (c *container) run(ev <-chan runtime.Event) {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	for {
		select {
		case <-c.exit:
			return
		case <-t.C:
			for _, pods := range c.pods {
				for _, pod := range pods {
					pod.restartIfDead(context.TODO())
				}
			}
		case _ = <-ev:
			// TODO: handle events
		}
	}
}

func (c *container) Stop() error {
	c.Lock()
	defer c.Unlock()

	select {
	case <-c.exit:
	default:
	}

	close(c.exit)
	for _, pods := range c.pods {
		for _, pod := range pods {
			// should we trace the error? but we are shutting down
			// the runtime, that's to say, the entire system is dying.
			pod.stop(context.TODO())
		}
	}
	c.wg.Wait()
	return nil
}

func New(c Client, opts ...runtime.Option) runtime.Runtime {
	var options runtime.Options
	for _, o := range opts {
		o(&options)
	}
	return &container{
		options: options,
		client:  c,
		pods:    make(map[string]map[string]*cpod),
		exit:    make(chan struct{}),
	}
}

type cpod struct {
	sync.RWMutex
	*runtime.Pod
	started int
	retries int
	running bool

	client    Client
	container *Container
	wg        *sync.WaitGroup
}

func (p *cpod) start(c context.Context) (err error) {
	p.Lock()
	defer p.Unlock()

	if !p.retry() {
		return
	}
	if err = p.client.Start(p.container, StartContext(c)); err != nil {
		return
	}
	p.running = true
	p.Status(runtime.Running, nil)
	p.Metadata["container_id"] = p.container.ID
	p.wg.Add(1)

	go p.wait(c)
	return nil
}

func (p *cpod) wait(c context.Context) {
	err := p.client.Wait(p.container, WaitContext(c))
	p.Lock()
	p.Status(runtime.Exited, err)
	p.running = false
	p.started++
	p.Unlock()
	p.wg.Done()
}

func (p *cpod) stop(c context.Context) error {
	p.Status(runtime.Stopping, nil)
	return p.client.Stop(p.container, StopContext(c))
}

func (p *cpod) retry() bool {
	if p.running {
		return false
	}
	return p.retries == -1 || p.started <= p.retries
}

func (p *cpod) restartIfDead(c context.Context) error {
	p.RLock()
	if !p.retry() {
		p.RUnlock()
		return nil
	}
	p.RUnlock()
	return p.start(c)
}

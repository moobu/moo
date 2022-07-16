package container

import (
	"context"
	"sync"
	"time"

	"github.com/moobu/moo/runtime"
)

// Client is an interface for managing containers.
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

	// let's first see if we already have this pod running
	id := pod.String()
	ns := options.Namespace
	if _, ok := c.pods[ns]; !ok {
		c.pods[ns] = make(map[string]*cpod)
	}
	if _, ok := c.pods[id]; ok {
		return runtime.ErrExists
	}

	// if not, we call our container runtime to create a container.
	ctx := options.Context
	con, err := c.client.Create(options.Bundle.Image, CreateContext(ctx))
	if err != nil {
		return err
	}

	// then create a pod wrapping this container and get the runtime
	// start the container.
	p := &cpod{
		Pod:       pod,
		wg:        &c.wg,
		container: con,
		retries:   options.Retries,
	}
	if err := p.start(ctx, c.client, false); err != nil {
		return err
	}
	// finally we save the pod in our memory.
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

	// let's first see if we have a pod wrapping the container.
	id := pod.String()
	ns := options.Namespace
	if _, ok := c.pods[ns]; !ok {
		return nil
	}
	p, ok := c.pods[ns][id]
	if !ok {
		return runtime.ErrNotFound
	}

	// then we have to see if the container really exists.
	ctx := options.Context
	con, err := c.client.Inspect(p.container.ID, InspectContext(ctx))
	if err != nil {
		return err
	}
	// we delete the container directly if it exited already.
	if con.Status == runtime.Exited {
		return c.client.Delete(con.ID, DeleteContext(ctx))
	}
	// otherwise we stop it before deleting it.
	if err := p.stop(ctx, c.client); err != nil {
		return err
	}
	// then we call our container runtime to remove this container
	if err := c.client.Delete(con.ID, DeleteContext(ctx)); err != nil {
		return err
	}
	// finally we remove the pod from our memory.
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

	// call the scheduler if we have one.
	var events <-chan runtime.Event
	if c.options.Scheduler != nil {
		if events, err = c.options.Scheduler.Schedule(); err != nil {
			return
		}
	}

	// start the runtime daemon.
	go c.run(events)
	return
}

func (c *container) run(events <-chan runtime.Event) {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	for {
		select {
		case <-c.exit:
			return
		case <-t.C:
			for _, pods := range c.pods {
				for _, pod := range pods {
					if err := pod.restartIfDead(context.TODO(), c.client); err != nil {
						// log error
					}
				}
			}
		case ev := <-events:
			// filter out orphaned or outdated events
			c.RLock()
			ns := ev.Namespace
			id := ev.Pod.String()
			pods, ok := c.pods[ns]
			if !ok {
				continue
			}
			pod, ok := pods[id]
			if !ok {
				continue
			}
			if ev.Time.Before(pod.updated) {
				continue
			}
			c.RUnlock()

			// handle the event
			if err := c.handle(&ev, pod); err != nil {
				// TODO: log error
			}
		}
	}
}

func (c *container) handle(ev *runtime.Event, pod *cpod) error {
	ctx := context.TODO()
	switch ev.Type {
	case runtime.EventStart:
		return pod.start(ctx, c.client, true) // force start
	case runtime.EventStop:
		return pod.stop(ctx, c.client)
	}
	return nil
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
			pod.stop(context.TODO(), c.client)
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
	updated time.Time

	container *Container
	wg        *sync.WaitGroup
}

func (p *cpod) start(ctx context.Context, c Client, force bool) (err error) {
	p.Lock()
	defer p.Unlock()

	if !force && !p.retry() {
		return
	}
	if err = c.Start(p.container, StartContext(ctx)); err != nil {
		return
	}
	p.running = true
	p.update(runtime.Running, nil)
	p.Metadata["container_id"] = p.container.ID
	p.wg.Add(1)

	go p.wait(ctx, c)
	return nil
}

func (p *cpod) update(status runtime.Status, err error) {
	now := time.Now()
	p.updated = now
	p.Status(status, now, err)
}

func (p *cpod) wait(ctx context.Context, c Client) {
	err := c.Wait(p.container, WaitContext(ctx))

	p.Lock()
	p.update(runtime.Exited, err)
	p.running = false
	p.started++
	p.Unlock()
	p.wg.Done()
}

func (p *cpod) stop(ctx context.Context, c Client) error {
	p.Lock()
	defer p.Unlock()

	p.update(runtime.Stopping, nil)
	return c.Stop(p.container, StopContext(ctx))
}

func (p *cpod) retry() bool {
	if p.running {
		return false
	}
	return p.retries == -1 || p.started <= p.retries
}

func (p *cpod) restartIfDead(ctx context.Context, c Client) error {
	p.RLock()
	if !p.retry() {
		p.RUnlock()
		return nil
	}
	p.RUnlock()
	return p.start(ctx, c, false)
}

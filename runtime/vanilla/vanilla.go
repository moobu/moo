package vanilla

import (
	"io"
	"sync"
	"time"

	"github.com/moobu/moo/builder"
	"github.com/moobu/moo/runtime"
)

// DEPRECATED

type Client interface {
	// Fork creates a new process group
	Fork(*Runnable) (*Process, error)
	// Kill kills the process and its child processes
	Kill(*Process) error
	// Wait waits the process to exit
	Wait(*Process) error
}

type Runnable struct {
	Bundle *builder.Bundle
	Env    []string
	Args   []string
}

type Process struct {
	ID     int
	Input  io.Writer
	Output io.Reader
	Error  io.Reader
}

type vanilla struct {
	sync.RWMutex
	options runtime.Options
	client  Client
	pods    map[string]map[string]*vpod
	wg      sync.WaitGroup
	exit    chan struct{}
}

func (v *vanilla) Create(pod *runtime.Pod, opts ...runtime.CreateOption) error {
	v.Lock()
	defer v.Unlock()

	var options runtime.CreateOptions
	for _, o := range opts {
		o(&options)
	}

	id := pod.String()
	ns := options.Namespace
	if _, ok := v.pods[ns]; !ok {
		v.pods[ns] = make(map[string]*vpod)
	}
	if _, ok := v.pods[ns][id]; ok {
		return runtime.ErrExists
	}

	vpod := &vpod{
		Pod:    pod,
		wg:     &v.wg,
		client: v.client,
		output: options.Output,
		runnable: &Runnable{
			Bundle: options.Bundle,
			Env:    options.Env,
			Args:   options.Args,
		},
	}

	if err := vpod.start(); err != nil {
		return err
	}
	v.pods[ns][id] = vpod
	return nil
}

func (v *vanilla) List(opts ...runtime.ListOption) ([]*runtime.Pod, error) {
	v.RLock()
	defer v.RUnlock()

	var options runtime.ListOptions
	for _, o := range opts {
		o(&options)
	}

	match := func(s, t string) bool {
		if len(t) == 0 {
			return true
		}
		return s == t
	}

	lpods, ok := v.pods[options.Namespace]
	if !ok {
		return nil, nil
	}

	pods := make([]*runtime.Pod, 0, len(lpods))
	for _, lpod := range lpods {
		if match(lpod.Name, options.Name) && match(lpod.Tag, options.Tag) {
			pods = append(pods, lpod.Pod)
		}
	}
	return pods, nil
}

func (v *vanilla) Delete(pod *runtime.Pod, opts ...runtime.DeleteOption) error {
	v.Lock()
	defer v.Unlock()

	var options runtime.DeleteOptions
	for _, o := range opts {
		o(&options)
	}

	key := pod.String()
	ns := options.Namespace
	if _, ok := v.pods[ns]; !ok {
		return nil
	}

	lpod, ok := v.pods[ns][key]
	if !ok {
		return runtime.ErrNotFound
	}
	if err := lpod.stop(); err != nil {
		return err
	}
	delete(v.pods[ns], key)
	return nil
}

func (v *vanilla) Start() error {
	go v.run() // start the runtime daemon
	return nil
}

func (v *vanilla) run() {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	for {
		select {
		case <-v.exit:
			return
		case <-t.C:
			v.RLock()
			for _, pods := range v.pods {
				for _, pod := range pods {
					pod.restartIfDead() // TODO: trace the error
				}
			}
			v.RUnlock()
		}
	}
}

func (v *vanilla) Stop() error {
	v.Lock()
	defer v.Unlock()

	select {
	case <-v.exit:
	default:
	}

	close(v.exit)
	for _, pods := range v.pods {
		for _, pod := range pods {
			// should we trace the error since we are shutting down
			// the runtime indicating the entire system is dying?
			pod.stop()
		}
	}
	v.wg.Wait()
	return nil
}

func New(driver Client, opts ...runtime.Option) runtime.Runtime {
	var options runtime.Options
	for _, o := range opts {
		o(&options)
	}
	return &vanilla{
		options: options,
		client:  driver,
		pods:    make(map[string]map[string]*vpod),
		exit:    make(chan struct{}),
	}
}

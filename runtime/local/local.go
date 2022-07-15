package local

import (
	"errors"
	"sync"
	"time"

	"github.com/moobu/moo/runtime"
	"github.com/moobu/moo/runtime/local/driver"
)

type local struct {
	sync.RWMutex
	options runtime.Options
	driver  driver.Driver
	pods    map[string]map[string]*lpod
	wg      sync.WaitGroup
	exit    chan struct{}
}

func (l *local) Create(pod *runtime.Pod, opts ...runtime.CreateOption) error {
	l.Lock()
	defer l.Unlock()

	var options runtime.CreateOptions
	for _, o := range opts {
		o(&options)
	}

	key := pod.String()
	ns := options.Namespace
	if _, ok := l.pods[ns]; !ok {
		l.pods[ns] = make(map[string]*lpod)
	}
	if _, ok := l.pods[ns][key]; ok {
		return errors.New("pod already created")
	}

	lpod := &lpod{
		Pod:    pod,
		wg:     &l.wg,
		driver: l.driver,
		output: options.Output,
		runnable: &driver.Runnable{
			Bundle: options.Bundle,
			Env:    options.Env,
			Args:   options.Args,
		},
	}

	if err := lpod.start(); err != nil {
		return err
	}
	l.pods[ns][key] = lpod
	return nil
}

func (l *local) List(opts ...runtime.ListOption) ([]*runtime.Pod, error) {
	l.RLock()
	defer l.RUnlock()

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

	lpods, ok := l.pods[options.Namespace]
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

func (l *local) Delete(pod *runtime.Pod, opts ...runtime.DeleteOption) error {
	l.Lock()
	defer l.Unlock()

	var options runtime.DeleteOptions
	for _, o := range opts {
		o(&options)
	}

	key := pod.String()
	ns := options.Namespace
	if _, ok := l.pods[ns]; !ok {
		return nil
	}

	lpod, ok := l.pods[ns][key]
	if !ok {
		return errors.New("no such pod")
	}
	if err := lpod.stop(); err != nil {
		return err
	}
	delete(l.pods[ns], key)
	return nil
}

func (l *local) Start() error {
	go l.run() // run the runtime daemon
	return nil
}

func (l *local) run() {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	for {
		select {
		case <-l.exit:
			return
		case <-t.C:
			l.RLock()
			for _, pods := range l.pods {
				for _, pod := range pods {
					pod.restartIfDead() // TODO: trace the error
				}
			}
			l.RUnlock()
		}
	}
}

func (l *local) Stop() error {
	l.Lock()
	defer l.Unlock()

	select {
	case <-l.exit:
	default:
	}

	close(l.exit)
	for _, pods := range l.pods {
		for _, pod := range pods {
			// should we trace the error since we are shutting down
			// the runtime indicating the entire system is dying?
			pod.stop()
		}
	}
	l.wg.Wait()
	return nil
}

func New(driver driver.Driver, opts ...runtime.Option) runtime.Runtime {
	var options runtime.Options
	for _, o := range opts {
		o(&options)
	}
	return &local{
		options: options,
		driver:  driver,
		pods:    make(map[string]map[string]*lpod),
		exit:    make(chan struct{}),
	}
}

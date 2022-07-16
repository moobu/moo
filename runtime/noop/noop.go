package noop

import (
	"errors"
	"sync"
	"time"

	"github.com/moobu/moo/runtime"
)

// This implementation does not really start or stop any pod, it
// just saves them so that we can test the components that need
// the runtime to maintain some pods but not the actual running ones.
type noop struct {
	sync.RWMutex
	options runtime.Options
	pods    map[string]map[string]*runtime.Pod
}

func (n *noop) Create(pod *runtime.Pod, opts ...runtime.CreateOption) error {
	n.Lock()
	defer n.Unlock()

	var options runtime.CreateOptions
	for _, o := range opts {
		o(&options)
	}

	key := pod.String()
	ns := options.Namespace
	if _, ok := n.pods[ns]; !ok {
		n.pods[ns] = make(map[string]*runtime.Pod)
	}
	if _, ok := n.pods[ns][key]; ok {
		return errors.New("pod already created")
	}

	pod.Status(runtime.Running, time.Now(), nil)
	pod.Metadata["started"] = time.Now().Format(time.RFC3339)
	pod.Metadata["source"] = options.Bundle.Source.URL
	n.pods[ns][key] = pod
	return nil
}

func (n *noop) Delete(pod *runtime.Pod, opts ...runtime.DeleteOption) error {
	n.Lock()
	defer n.Unlock()

	var options runtime.DeleteOptions
	for _, o := range opts {
		o(&options)
	}

	key := pod.String()
	ns := options.Namespace
	if _, ok := n.pods[ns]; !ok {
		return errors.New("no such namespace")
	}

	if _, ok := n.pods[ns][key]; !ok {
		return errors.New("no such pod")
	}

	delete(n.pods[ns], key)
	return nil
}

func (n *noop) List(opts ...runtime.ListOption) ([]*runtime.Pod, error) {
	n.RLock()
	defer n.RUnlock()

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

	saved, ok := n.pods[options.Namespace]
	if !ok {
		return nil, nil
	}

	pods := make([]*runtime.Pod, 0, len(saved))
	for _, lpod := range saved {
		if match(lpod.Name, options.Name) && match(lpod.Tag, options.Tag) {
			pods = append(pods, lpod)
		}
	}
	return pods, nil
}

func (n *noop) Start() error {
	return nil
}

func (n *noop) Stop() error {
	return nil
}

func New(opts ...runtime.Option) runtime.Runtime {
	var options runtime.Options
	for _, o := range opts {
		o(&options)
	}
	return &noop{
		options: options,
		pods:    make(map[string]map[string]*runtime.Pod),
	}
}

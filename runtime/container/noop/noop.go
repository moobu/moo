package noop

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sync"

	"github.com/moobu/moo/runtime"
	"github.com/moobu/moo/runtime/container"
)

type con struct {
	id     string
	image  string
	status runtime.Status
	exit   chan struct{}
}

type noop struct {
	sync.RWMutex
	options    container.Options
	containers map[string]*con
	nextid     int32
}

func (n *noop) Create(image string, opts ...container.CreateOption) (*container.Container, error) {
	n.Lock()
	defer n.Unlock()

	id := r6()
	n.containers[id] = &con{
		id:    id,
		image: image,
	}
	return &container.Container{
		ID:     id,
		Image:  image,
		Status: runtime.Pending,
	}, nil
}

func (n *noop) Inspect(id string, opts ...container.InspectOption) (*container.Container, error) {
	n.RLock()
	defer n.RUnlock()

	con, ok := n.containers[id]
	if !ok {
		return nil, errors.New("no such container")
	}
	return &container.Container{
		ID:     id,
		Image:  con.image,
		Status: con.status,
	}, nil
}

func (n *noop) Delete(id string, opts ...container.DeleteOption) error {
	n.Lock()
	defer n.Unlock()

	var options container.DeleteOptions
	for _, o := range opts {
		o(&options)
	}

	con, ok := n.containers[id]
	if !ok {
		return errors.New("no such container")
	}
	if con.status == runtime.Running && !options.Force {
		return errors.New("the container is still running")
	}
	delete(n.containers, con.id)
	return nil
}

func (n *noop) List(...container.ListOption) ([]*container.Container, error) {
	n.RLock()
	defer n.RUnlock()

	cons := make([]*container.Container, len(n.containers))
	for k, v := range n.containers {
		cons = append(cons, &container.Container{
			ID:     k,
			Image:  v.image,
			Status: v.status,
		})
	}
	return cons, nil
}

func (n *noop) Wait(c *container.Container, opts ...container.WaitOption) error {
	n.RLock()

	con, ok := n.containers[c.ID]
	if !ok {
		return errors.New("no such container")
	}
	n.RUnlock()
	<-con.exit
	return nil
}

func (n *noop) Start(c *container.Container, opts ...container.StartOption) error {
	n.Lock()
	defer n.Unlock()

	con, ok := n.containers[c.ID]
	if !ok {
		return errors.New("no such container")
	}
	con.exit = make(chan struct{})
	con.status = runtime.Running
	return nil
}

func (n *noop) Stop(c *container.Container, opts ...container.StopOption) error {
	n.Lock()
	defer n.Unlock()

	con, ok := n.containers[c.ID]
	if !ok {
		return errors.New("no such container")
	}

	select {
	case <-con.exit:
	default:
	}
	close(con.exit)
	con.status = runtime.Exited
	return nil
}

func New(opts ...container.Option) container.Client {
	var options container.Options
	for _, o := range opts {
		o(&options)
	}
	return &noop{
		options:    options,
		containers: make(map[string]*con),
	}
}

func r6() string {
	var b [6]byte
	rand.Read(b[:])
	return fmt.Sprintf("%x", b)
}

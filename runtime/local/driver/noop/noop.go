package noop

import (
	"bytes"
	"errors"
	"io"
	"sync"

	"github.com/moobu/moo/runtime/local/driver"
)

type proc struct {
	id   int32
	exit chan struct{}
}

type noop struct {
	sync.Mutex
	procs   map[int32]*proc
	nextpid int32
}

func (n *noop) Fork(r *driver.Runnable) (*driver.Process, error) {
	n.Lock()
	defer n.Unlock()
	n.nextpid++
	pid := n.nextpid
	n.procs[pid] = &proc{
		id:   pid,
		exit: make(chan struct{}),
	}
	return &driver.Process{
		ID:  int(pid),
		In:  io.Discard,
		Out: new(bytes.Buffer),
		Err: new(bytes.Buffer),
	}, nil
}

func (n *noop) Kill(p *driver.Process) error {
	n.Lock()
	defer n.Unlock()

	proc, ok := n.procs[int32(p.ID)]
	if !ok {
		return errors.New("no such process")
	}
	close(proc.exit)
	delete(n.procs, proc.id)
	return nil
}

func (n *noop) Wait(p *driver.Process) error {
	n.Lock()
	proc, ok := n.procs[int32(p.ID)]
	if !ok {
		n.Unlock()
		return errors.New("no such process")
	}
	n.Unlock()
	<-proc.exit
	return nil
}

func New() driver.Driver {
	return &noop{}
}

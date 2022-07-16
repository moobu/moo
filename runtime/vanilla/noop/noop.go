package noop

import (
	"bytes"
	"errors"
	"io"
	"sync"

	"github.com/moobu/moo/runtime/vanilla"
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

func (n *noop) Fork(r *vanilla.Runnable) (*vanilla.Process, error) {
	n.Lock()
	defer n.Unlock()
	n.nextpid++
	pid := n.nextpid
	n.procs[pid] = &proc{
		id:   pid,
		exit: make(chan struct{}),
	}
	return &vanilla.Process{
		ID:     int(pid),
		Input:  io.Discard,
		Output: new(bytes.Buffer),
		Error:  new(bytes.Buffer),
	}, nil
}

func (n *noop) Kill(p *vanilla.Process) error {
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

func (n *noop) Wait(p *vanilla.Process) error {
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

func New() vanilla.Client {
	return &noop{}
}

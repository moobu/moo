package vanilla

import (
	"io"
	"sync"

	"github.com/moobu/moo/runtime"
)

type vpod struct {
	sync.RWMutex
	*runtime.Pod
	started int // retries the pod has consumed
	retries int
	running bool // avoid exceedingly restarting the pod

	output   io.Writer
	client   Client
	runnable *Runnable
	process  *Process
	wg       *sync.WaitGroup
}

func (p *vpod) restartIfDead() error {
	p.RLock()
	if !p.retry() {
		p.RUnlock()
		return nil
	}
	p.RUnlock()
	return p.start()
}

func (p *vpod) retry() bool {
	if p.running {
		return false
	}
	return p.retries == -1 || p.started <= p.retries
}

func (p *vpod) start() (err error) {
	p.Lock()
	defer p.Unlock()

	if !p.retry() {
		return
	}

	p.process, err = p.client.Fork(p.runnable)
	if err != nil {
		return
	}

	p.Status(runtime.Running, nil)
	p.running = true
	p.wg.Add(1)

	if p.output != nil {
		p.stream()
	}
	go p.wait()
	return nil
}

func (p *vpod) stream() {
	go io.Copy(p.output, p.process.Output)
	go io.Copy(p.output, p.process.Error)
}

func (p *vpod) wait() {
	err := p.client.Wait(p.process)
	p.Lock()
	p.Status(runtime.Exited, err)
	p.started++
	p.running = false
	p.wg.Done()
	p.Unlock()
}

func (p *vpod) stop() error {
	p.Status(runtime.Stopping, nil)
	return p.client.Kill(p.process)
}

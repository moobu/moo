package runtime

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrExists   = errors.New("pod already exists")
	ErrNotFound = errors.New("no such pod")
)

type Runtime interface {
	// Create creates a pod containing a process
	Create(*Pod, ...CreateOption) error
	// Delete deletes a pod
	Delete(*Pod, ...DeleteOption) error
	// List lists pods
	List(...ListOption) ([]*Pod, error)
	// Start starts the runtime
	Start() error
	// Start stops the runtime
	Stop() error
}

type Pod struct {
	Name     string
	Tag      string
	Metadata map[string]string
}

func (p Pod) String() string {
	return fmt.Sprintf("%s:%s", p.Name, p.Tag)
}

func Parse(rawPod string) (*Pod, error) {
	pod := &Pod{}
	i := strings.IndexByte(rawPod, ':')
	if i != -1 {
		pod.Tag = rawPod[i+1:]
	}
	pod.Name = rawPod[:i]
	return pod, nil
}

func (p *Pod) Status(status Status, err error) {
	if p.Metadata == nil {
		p.Metadata = make(map[string]string)
	}
	if err != nil {
		p.Metadata["error"] = err.Error()
	}
	p.Metadata["status"] = status.String()
	p.Metadata["updated"] = time.Now().Format(time.RFC3339)
}

func (p Pod) Get(key string) string {
	if p.Metadata == nil {
		return "N/A"
	}
	if v, ok := p.Metadata[key]; ok {
		return v
	}
	return "N/A"
}

type Status int8

func (s Status) String() string {
	return StatusText[s]
}

const (
	Pending Status = iota
	Running
	Stopping
	Exited
)

var StatusText = [...]string{
	"pendding",
	"running",
	"stopping",
	"exited",
}

type Scheduler interface {
	Schedule() (<-chan Event, error)
}

type Event struct {
	Type EventType
	Time time.Time
	Pod  *Pod
}

type EventType int8

const (
	EventUpdate EventType = iota + 1
	EventStart
	EventStop
)

var Default Runtime

func Create(pod *Pod, opts ...CreateOption) error {
	return Default.Create(pod, opts...)
}

func Delete(pod *Pod, opts ...DeleteOption) error {
	return Default.Delete(pod, opts...)
}

func List(opts ...ListOption) ([]*Pod, error) {
	return Default.List(opts...)
}

func Start() error {
	return Default.Start()
}

func Stop() error {
	return Default.Stop()
}

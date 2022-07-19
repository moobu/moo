package router

import (
	"errors"
	"hash/fnv"
)

var (
	ErrDuplicated = errors.New("duplicated route")
	ErrNotFound   = errors.New("no such route")
)

type Router interface {
	// Register adds a route to a pod
	Register(*Route) error
	// Deregister removes a route to a pod
	Deregister(*Route) error
	// Lookup finds all routes to the same pod
	Lookup(string) ([]*Route, error)
}

type Route struct {
	Path     string // pod path (fmt. /:namespace/:name/:tag)
	Protocol string // the protocol by which we communicate with the pod
	Address  string // pod address (e.g. 10.0.0.1:80, /tmp/moo/xxx.sock)
}

func (r Route) Sum() uint32 {
	h := fnv.New32()
	h.Write([]byte(r.Path + r.Protocol + r.Address))
	return h.Sum32()
}

var Default Router

func Register(r *Route) error {
	return Default.Register(r)
}

func Deregister(r *Route) error {
	return Default.Deregister(r)
}

func Lookup(pod string) ([]*Route, error) {
	return Default.Lookup(pod)
}

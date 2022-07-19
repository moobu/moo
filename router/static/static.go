package static

import (
	"sync"

	"github.com/moobu/moo/router"
)

type static struct {
	sync.RWMutex
	options router.Options
	routes  map[string]map[uint32]*router.Route
}

func (s *static) Register(r *router.Route) error {
	s.Lock()
	defer s.Unlock()

	pod := r.Path
	if _, ok := s.routes[r.Path]; !ok {
		s.routes[pod] = make(map[uint32]*router.Route)
	}
	sum := r.Sum()
	if _, ok := s.routes[pod][sum]; ok {
		return router.ErrDuplicated
	}
	s.routes[pod][sum] = r
	return nil
}

func (s *static) Deregister(r *router.Route) error {
	s.Lock()
	defer s.Unlock()

	pod := r.Path
	if _, ok := s.routes[pod]; !ok {
		return nil
	}
	sum := r.Sum()
	if _, ok := s.routes[pod][sum]; !ok {
		return nil
	}
	delete(s.routes[pod], sum)
	return nil
}

func (s *static) Lookup(pod string) ([]*router.Route, error) {
	s.RLock()
	defer s.RUnlock()

	routes, ok := s.routes[pod]
	if !ok {
		return nil, router.ErrNotFound
	}

	// clone in case that routes change
	clone := make([]*router.Route, 0, len(routes))
	for _, route := range routes {
		clone = append(clone, route)
	}
	return clone, nil
}

func New(opts ...router.Option) router.Router {
	var options router.Options
	for _, o := range opts {
		o(&options)
	}
	return &static{
		options: options,
		routes:  make(map[string]map[uint32]*router.Route),
	}
}

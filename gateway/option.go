package gateway

import "github.com/moobu/moo/router"

type Options struct {
}

type Option func(*Options)

type ProxyOptions struct {
	Router router.Router
}

type ProxyOption func(*ProxyOptions)

func Router(r router.Router) ProxyOption {
	return func(o *ProxyOptions) {
		o.Router = r
	}
}

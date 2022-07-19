package gateway

import "github.com/moobu/moo/router"

type Options struct {
	Domain string
}

type Option func(*Options)

func Domain(name string) Option {
	return func(o *Options) {
		o.Domain = name
	}
}

type ProxyOptions struct {
	Router router.Router
}

type ProxyOption func(*ProxyOptions)

func Router(r router.Router) ProxyOption {
	return func(o *ProxyOptions) {
		o.Router = r
	}
}

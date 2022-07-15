package client

type Options struct {
	Server string
}

type Option func(*Options)

func Server(addr string) Option {
	return func(o *Options) {
		o.Server = addr
	}
}

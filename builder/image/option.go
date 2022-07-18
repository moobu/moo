package image

import "context"

type Options struct {
	Addr string
}

type Option func(*Options)

func Addr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

type BuildOptions struct {
	Context context.Context
	Name    string
	Tag     string
}

type BuildOption func(*BuildOptions)

func BuildContext(c context.Context) BuildOption {
	return func(o *BuildOptions) {
		o.Context = c
	}
}

func Name(name string) BuildOption {
	return func(o *BuildOptions) {
		o.Name = name
	}
}

func Tag(tag string) BuildOption {
	return func(o *BuildOptions) {
		o.Tag = tag
	}
}

type RemoveOptions struct {
	Context context.Context
}

type RemoveOption func(*RemoveOptions)

func RemoveContext(c context.Context) RemoveOption {
	return func(o *RemoveOptions) {
		o.Context = c
	}
}

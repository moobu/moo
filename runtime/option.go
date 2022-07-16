package runtime

import (
	"context"
	"io"

	"github.com/moobu/moo/builder"
)

type Options struct {
	Scheduler Scheduler
}

type Option func(*Options)

func WithScheduler(s Scheduler) Option {
	return func(o *Options) {
		o.Scheduler = s
	}
}

type CreateOptions struct {
	Context   context.Context
	Bundle    *builder.Bundle
	Output    io.Writer `json:"-"`
	Env       []string
	Args      []string
	Namespace string
	Replicas  int
	Retries   int
	GPU       bool
}

type CreateOption func(*CreateOptions)

func CreateContext(c context.Context) CreateOption {
	return func(o *CreateOptions) {
		o.Context = c
	}
}

func Env(env ...string) CreateOption {
	return func(o *CreateOptions) {
		o.Env = env
	}
}

func Args(args ...string) CreateOption {
	return func(o *CreateOptions) {
		o.Args = args
	}
}

func Bundle(bundle *builder.Bundle) CreateOption {
	return func(o *CreateOptions) {
		o.Bundle = bundle
	}
}

func Output(w io.Writer) CreateOption {
	return func(o *CreateOptions) {
		o.Output = w
	}
}

func GPU(enable bool) CreateOption {
	return func(o *CreateOptions) {
		o.GPU = enable
	}
}

func CreateNamespace(ns string) CreateOption {
	return func(o *CreateOptions) {
		o.Namespace = ns
	}
}

func Replicas(n int) CreateOption {
	return func(o *CreateOptions) {
		o.Replicas = n
	}
}

func Retries(n int) CreateOption {
	return func(o *CreateOptions) {
		o.Retries = n
	}
}

type ListOptions struct {
	Context   context.Context
	Namespace string
	Name      string
	Tag       string
	Verbose   bool
}

type ListOption func(*ListOptions)

func ListContext(c context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = c
	}
}

func ListNamespace(ns string) ListOption {
	return func(o *ListOptions) {
		o.Namespace = ns
	}
}

func Name(name string) ListOption {
	return func(o *ListOptions) {
		o.Name = name
	}
}

func Tag(tag string) ListOption {
	return func(o *ListOptions) {
		o.Tag = tag
	}
}

func Verbose(v bool) ListOption {
	return func(o *ListOptions) {
		o.Verbose = v
	}
}

type DeleteOptions struct {
	Context   context.Context
	Namespace string
}

type DeleteOption func(*DeleteOptions)

func DeleteContext(c context.Context) DeleteOption {
	return func(o *DeleteOptions) {
		o.Context = c
	}
}

func DeleteNamespace(ns string) DeleteOption {
	return func(o *DeleteOptions) {
		o.Namespace = ns
	}
}

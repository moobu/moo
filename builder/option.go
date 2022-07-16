package builder

import "context"

type Options struct{}

type Option func(*Options)

type BuildOptions struct {
	Context   context.Context
	Namespace string
	Dir       string
	Ref       string
}

type BuildOption func(*BuildOptions)

func BuildContext(c context.Context) BuildOption {
	return func(o *BuildOptions) {
		o.Context = c
	}
}

func BuildNamespace(ns string) BuildOption {
	return func(o *BuildOptions) {
		o.Namespace = ns
	}
}

func Dir(dir string) BuildOption {
	return func(o *BuildOptions) {
		o.Dir = dir
	}
}

func Ref(ref string) BuildOption {
	return func(o *BuildOptions) {
		o.Ref = ref
	}
}

type CleanOptions struct {
	Context   context.Context
	Namespace string
}

type CleanOption func(*CleanOptions)

func CleanContext(c context.Context) CleanOption {
	return func(o *CleanOptions) {
		o.Context = c
	}
}

func CleanNamespace(ns string) CleanOption {
	return func(o *CleanOptions) {
		o.Namespace = ns
	}
}

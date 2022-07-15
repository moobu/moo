package builder

type Options struct{}

type Option func(*Options)

type BuildOptions struct {
	Namespace string
	Dir       string
	Ref       string
}

type BuildOption func(*BuildOptions)

func BuildWithNamespace(ns string) BuildOption {
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
	Namespace string
}

type CleanOption func(*CleanOptions)

func CleanWithNamespace(ns string) CleanOption {
	return func(o *CleanOptions) {
		o.Namespace = ns
	}
}

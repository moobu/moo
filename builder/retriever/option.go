package retriever

type Options struct {
	Dir string
}

type Option func(*Options)

func Dir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

type RetrieveOptions struct {
	Ref string
}

type RetrieveOption func(*RetrieveOptions)

func Ref(ref string) RetrieveOption {
	return func(o *RetrieveOptions) {
		o.Ref = ref
	}
}

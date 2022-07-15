package noop

import (
	"github.com/moobu/moo/builder"
)

type noop struct {
	options builder.Options
}

func (n noop) Build(s *builder.Source, opts ...builder.BuildOption) (*builder.Bundle, error) {
	var options builder.BuildOptions
	for _, o := range opts {
		o(&options)
	}
	return &builder.Bundle{
		Dir:    options.Dir,
		Ref:    options.Ref,
		Source: s,
	}, nil
}

func (n noop) Clean(b *builder.Bundle, opts ...builder.CleanOption) error {
	return nil
}

func (noop) String() string {
	return "noop"
}

func New(opts ...builder.Option) builder.Builder {
	var options builder.Options
	for _, o := range opts {
		o(&options)
	}
	return &noop{options: options}
}

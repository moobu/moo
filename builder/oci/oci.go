package podman

import "github.com/moobu/moo/builder"

// TODO: implement

type oci struct {
	options builder.Options
}

func (o *oci) Build(s *builder.Source, opts ...builder.BuildOption) (*builder.Bundle, error) {
	return nil, nil
}

func (o *oci) Clean(b *builder.Bundle, opts ...builder.CleanOption) error {
	return nil
}

func (o oci) String() string {
	return "oci"
}

func New(opts ...builder.Option) builder.Builder {
	var options builder.Options
	for _, o := range opts {
		o(&options)
	}
	return &oci{options: options}
}

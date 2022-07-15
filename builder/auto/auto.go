package auto

import (
	"fmt"

	"github.com/moobu/moo/builder"
)

type auto struct {
	options builder.Options
	langs   map[string]builder.Builder
}

func (a auto) Build(s *builder.Source, opts ...builder.BuildOption) (*builder.Bundle, error) {
	if lang, ok := a.langs[s.Type]; ok {
		return lang.Build(s, opts...)
	}
	return nil, fmt.Errorf("no builder implemented for %s", s.Type)
}

func (a auto) Clean(b *builder.Bundle, opts ...builder.CleanOption) error {
	if lang, ok := a.langs[b.Source.Type]; ok {
		return lang.Clean(b, opts...)
	}
	return fmt.Errorf("no builder implemented for %s", b.Source.Type)
}

func (auto) String() string {
	return "auto"
}

// New returns a auto builder that automatically specifies which named
// builder to use, according to the type of the source.
func New(bs []builder.Builder, opts ...builder.Option) builder.Builder {
	var options builder.Options
	for _, o := range opts {
		o(&options)
	}

	langs := make(map[string]builder.Builder, len(bs))
	for _, b := range bs {
		langs[b.String()] = b
	}
	return &auto{
		options: options,
		langs:   langs,
	}
}

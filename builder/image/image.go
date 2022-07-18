package image

import (
	"github.com/moobu/moo/builder"
)

// Client is an interface for managing images.
type Client interface {
	Build(string, ...BuildOption) (*Image, error)
	Remove(string, ...RemoveOption) error
}

type Image struct {
	ID string
}

type image struct {
	options builder.Options
	client  Client
}

func (i *image) Build(s *builder.Source, opts ...builder.BuildOption) (*builder.Bundle, error) {
	var options builder.BuildOptions
	for _, o := range opts {
		o(&options)
	}

	im, err := i.client.Build(options.Dir, BuildContext(options.Context), Name(s.Name), Tag(options.Ref))
	if err != nil {
		return nil, err
	}
	return &builder.Bundle{
		Source: s,
		Image:  im.ID,
		Dir:    options.Dir,
		Ref:    options.Ref,
	}, nil
}

func (i *image) Clean(b *builder.Bundle, opts ...builder.CleanOption) error {
	var options builder.CleanOptions
	for _, o := range opts {
		o(&options)
	}
	return i.client.Remove(b.Image, RemoveContext(options.Context))
}

func (image) String() string {
	return "image"
}

func New(c Client, opts ...builder.Option) builder.Builder {
	var options builder.Options
	for _, o := range opts {
		o(&options)
	}
	return &image{
		options: options,
		client:  c,
	}
}

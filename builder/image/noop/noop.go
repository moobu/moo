package noop

import (
	"crypto/rand"
	"fmt"

	"github.com/moobu/moo/builder"
	"github.com/moobu/moo/builder/image"
)

type noop struct {
	options builder.Options
}

func (n *noop) Build(path string, opts ...image.BuildOption) (*image.Image, error) {
	return &image.Image{
		ID: r6(),
	}, nil
}

func (n *noop) Remove(id string, opts ...image.RemoveOption) error {
	return nil
}

func New(opts ...image.Option) image.Client {
	return &noop{}
}

func r6() string {
	var b [6]byte
	rand.Read(b[:])
	return fmt.Sprintf("%x", b)
}

package http

import (
	"encoding/json"
	"fmt"

	"github.com/moobu/moo/builder"
	"github.com/moobu/moo/internal/pool/buffer"
	server "github.com/moobu/moo/server/http"
)

func (h *http) Build(s *builder.Source, opts ...builder.BuildOption) (*builder.Bundle, error) {
	var options builder.BuildOptions
	for _, o := range opts {
		o(&options)
	}

	url := fmt.Sprintf("http://%s/build", h.options.Server)
	reader := buffer.Get()
	defer buffer.Put(reader)

	args := server.BuildArgs{Source: s, Options: &options}
	enc := json.NewEncoder(reader)
	if err := enc.Encode(args); err != nil {
		return nil, err
	}

	body, err := h.invoke("POST", url, reader)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	bundle := &builder.Bundle{}
	dec := json.NewDecoder(body)
	if err := dec.Decode(bundle); err != nil {
		return nil, err
	}
	return bundle, nil
}

func (h *http) Clean(b *builder.Bundle, opts ...builder.CleanOption) error {
	var options builder.CleanOptions
	for _, o := range opts {
		o(&options)
	}

	args := server.CleanArgs{Bundle: b, Options: &options}
	url := fmt.Sprintf("http://%s/clean", h.options.Server)

	reader := buffer.Get()
	defer buffer.Put(reader)

	encoder := json.NewEncoder(reader)
	if err := encoder.Encode(args); err != nil {
		return err
	}
	_, err := h.invoke("POST", url, reader)
	return err
}

package http

import (
	"encoding/json"
	"fmt"

	"github.com/moobu/moo/internal/pool/buffer"
	"github.com/moobu/moo/router"
)

func (h *http) Register(route *router.Route) error {
	url := fmt.Sprintf("http://%s/register", h.options.Server)

	reader := buffer.Get()
	defer buffer.Put(reader)

	encoder := json.NewEncoder(reader)
	if err := encoder.Encode(route); err != nil {
		return err
	}
	_, err := h.invoke("POST", url, reader)
	return err
}

func (h *http) Deregister(route *router.Route) error {
	url := fmt.Sprintf("http://%s/deregister", h.options.Server)

	reader := buffer.Get()
	defer buffer.Put(reader)

	encoder := json.NewEncoder(reader)
	if err := encoder.Encode(route); err != nil {
		return err
	}
	_, err := h.invoke("POST", url, reader)
	return err
}

func (h *http) Lookup(pod string) ([]*router.Route, error) {
	route := router.Route{Pod: pod}
	url := fmt.Sprintf("http://%s/lookup", h.options.Server)

	reader := buffer.Get()
	defer buffer.Put(reader)

	encoder := json.NewEncoder(reader)
	if err := encoder.Encode(route); err != nil {
		return nil, err
	}
	body, err := h.invoke("POST", url, reader)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var routes []*router.Route
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&route); err != nil {
		return nil, err
	}
	return routes, nil
}

package http

import (
	"encoding/json"
	"fmt"

	"github.com/moobu/moo/internal/pool/buffer"
	"github.com/moobu/moo/runtime"
	server "github.com/moobu/moo/server/http"
)

func (h *http) Create(pod *runtime.Pod, opts ...runtime.CreateOption) error {
	var options runtime.CreateOptions
	for _, o := range opts {
		o(&options)
	}

	url := fmt.Sprintf("http://%s/create", h.options.Server)
	reader := buffer.Get()
	defer buffer.Put(reader)

	encoder := json.NewEncoder(reader)
	args := server.CreateArgs{Pod: pod, Options: &options}
	if err := encoder.Encode(args); err != nil {
		return err
	}
	_, err := h.invoke("POST", url, reader)
	return err
}

func (h *http) Delete(pod *runtime.Pod, opts ...runtime.DeleteOption) error {
	var options runtime.DeleteOptions
	for _, o := range opts {
		o(&options)
	}

	url := fmt.Sprintf("http://%s/delete", h.options.Server)
	reader := buffer.Get()
	defer buffer.Put(reader)

	encoder := json.NewEncoder(reader)
	args := server.DeleteArgs{Pod: pod, Options: &options}
	if err := encoder.Encode(args); err != nil {
		return err
	}
	_, err := h.invoke("POST", url, reader)
	return err
}

func (h *http) List(opts ...runtime.ListOption) ([]*runtime.Pod, error) {
	var options runtime.ListOptions
	for _, o := range opts {
		o(&options)
	}

	url := fmt.Sprintf("http://%s/list", h.options.Server)
	reader := buffer.Get()
	defer buffer.Put(reader)

	encoder := json.NewEncoder(reader)
	if err := encoder.Encode(&options); err != nil {
		return nil, err
	}

	body, err := h.invoke("POST", url, reader)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var pods []*runtime.Pod
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&pods); err != nil {
		return nil, err
	}
	return pods, nil
}

// the client implementation does not need the two methods below
// so we just leave them empty.
func (h *http) Start() error {
	return nil
}

func (h *http) Stop() error {
	return nil
}

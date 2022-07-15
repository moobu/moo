package http

import (
	"errors"
	"fmt"
	"io"
	std "net/http"

	"github.com/moobu/moo/client"
)

type http struct {
	options client.Options
	client  *std.Client
}

func New(opts ...client.Option) client.Client {
	var options client.Options
	for _, o := range opts {
		o(&options)
	}
	return &http{
		options: options,
		client:  &std.Client{},
	}
}

func (http) String() string {
	return "client"
}

// this is when HTTP sucks!
func (h *http) invoke(method, url string, body io.Reader) (io.ReadCloser, error) {
	req, err := std.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != std.StatusOK {
		return nil, fmt.Errorf("moo server responded with status %d", res.StatusCode)
	}
	merr := res.Header.Get("X-Moo-Error")
	if len(merr) > 0 {
		return nil, errors.New(merr)
	}
	return res.Body, nil
}

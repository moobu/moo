package auto

import (
	"fmt"
	"net/url"

	"github.com/moobu/moo/builder/retriever"
)

type auto struct {
	options retriever.Options
	protos  map[string]retriever.Retriever
}

func (a *auto) Retrieve(rawURL string, opts ...retriever.RetrieveOption) (*retriever.Repository, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if proto, ok := a.protos[u.Scheme]; ok {
		return proto.Retrieve(rawURL, opts...)
	}
	return nil, fmt.Errorf("no retriever implemented for %s", u.Scheme)
}

func (a auto) String() string {
	return "auto"
}

// New returns an auto retriever that automatically chooses
// a protocol-specific retriever by the scheme of the URL.
func New(rs []retriever.Retriever, opts ...retriever.Option) retriever.Retriever {
	var options retriever.Options
	for _, o := range opts {
		o(&options)
	}
	protos := make(map[string]retriever.Retriever, len(rs))
	for _, r := range rs {
		protos[r.String()] = r
	}
	return &auto{
		options: options,
		protos:  protos,
	}
}

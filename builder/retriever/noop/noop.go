package noop

import (
	"net/url"

	"github.com/moobu/moo/builder/retriever"
)

type noop struct{}

func (n noop) Retrieve(rawURL string, opts ...retriever.RetrieveOption) (*retriever.Repository, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return &retriever.Repository{Path: u.Path}, nil
}

func (n noop) String() string {
	return "noop"
}

func New(opts ...retriever.Option) retriever.Retriever {
	return noop{}
}

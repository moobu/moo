package git

import "github.com/moobu/moo/builder/retriever"

// TODO: implement the git retriever
type git struct {
	options retriever.Options
}

func (g *git) Retrieve(url string, opts ...retriever.RetrieveOption) (*retriever.Repository, error) {
	var options retriever.RetrieveOptions
	for _, o := range opts {
		o(&options)
	}

	return nil, nil
}

func (g git) String() string {
	return "git"
}

func New(opts ...retriever.Option) retriever.Retriever {
	var options retriever.Options
	for _, o := range opts {
		o(&options)
	}
	return &git{options: options}
}

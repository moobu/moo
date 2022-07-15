package retriever

import (
	"errors"

	"github.com/moobu/moo/builder"
)

type Retriever interface {
	// Retrieve retrieves a source typically a git repository
	// from remote platforms like GitHub.
	Retrieve(string, ...RetrieveOption) (*Repository, error)
	// String returns the name of the implementation
	String() string
}

type Repository struct {
	Path string
}

type retriever struct {
	options builder.Options
	next    builder.Builder
	r       Retriever
}

func (r *retriever) Build(s *builder.Source, opts ...builder.BuildOption) (*builder.Bundle, error) {
	if len(s.URL) == 0 {
		return nil, errors.New("missing the source address")
	}
	repo, err := r.r.Retrieve(s.URL)
	if err != nil {
		return nil, err
	}
	// rewrite Dir
	opts = append(opts, builder.Dir(repo.Path)) // this is weired but I like it :sweating:
	return r.next.Build(s, opts...)
}

func (r retriever) Clean(bun *builder.Bundle, opts ...builder.CleanOption) error {
	return r.next.Clean(bun, opts...)
}

func (r retriever) String() string {
	return "retriever"
}

// New creates a retriever builder whose Build method retrieves the source
// code from the given URL of the remote repository and then uses the next
// builder to build the source code.
func New(r Retriever, next builder.Builder, opts ...builder.Option) builder.Builder {
	var options builder.Options
	for _, o := range opts {
		o(&options)
	}
	return &retriever{
		options: options,
		next:    next,
		r:       r,
	}
}

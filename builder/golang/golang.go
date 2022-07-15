package golang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/moobu/moo/builder"
)

type golang struct {
	options builder.Options
}

func (g *golang) Build(s *builder.Source, opts ...builder.BuildOption) (*builder.Bundle, error) {
	var options builder.BuildOptions
	for _, o := range opts {
		o(&options)
	}

	bin := s.Name
	cmd := exec.Command("go", "build", "-o", bin)
	// TODO: use a session logger to write the output to
	cmd.Stdout = os.Stdout
	cmd.Dir = options.Dir
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return &builder.Bundle{
		Ref:    options.Ref,
		Dir:    options.Dir,
		Entry:  []string{fmt.Sprintf("./%s", bin)},
		Source: s,
	}, nil
}

func (g *golang) Clean(b *builder.Bundle, opts ...builder.CleanOption) error {
	// we clean up just by removing the binary
	return os.Remove(filepath.Join(b.Dir, b.Source.Name))
}

func (g golang) String() string {
	return "go"
}

// New returns a python builder that uses conda for dependencies isolation.
func New(opts ...builder.Option) builder.Builder {
	var options builder.Options
	for _, o := range opts {
		o(&options)
	}
	return &golang{options: options}
}

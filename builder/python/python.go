package python

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/moobu/moo/builder"
)

type python struct {
	options builder.Options
}

func (p *python) Build(s *builder.Source, opts ...builder.BuildOption) (*builder.Bundle, error) {
	var options builder.BuildOptions
	for _, o := range opts {
		o(&options)
	}
	// we first create a new conda environment for the source
	env := fmt.Sprintf("%s-%s", options.Namespace, s.Name)
	cmd := exec.Command("conda", "create", "-y", "-n", env)
	// TODO: use a session logger to write the output to
	cmd.Stdout = os.Stdout
	cmd.Dir = options.Dir
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	requires := "requirements.txt"
	if _, err := os.Stat(requires); err != nil {
		return nil, errors.New("requirements.txt is not found or invalid")
	}
	// then we install the dependencies to the newly created environment
	cmd = exec.Command("conda", "install", "-r", requires, "-y")
	// TODO: use a session logger to write the output to
	cmd.Stdout = os.Stdout
	cmd.Dir = options.Dir
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	// return the bundle whose Entry provides the pre-command for the
	// runtime to execute.
	return &builder.Bundle{
		Dir:    options.Dir,
		Ref:    options.Ref,
		Entry:  []string{"conda", "run", "-n", env},
		Source: s,
	}, nil
}

func (p *python) Clean(b *builder.Bundle, opts ...builder.CleanOption) error {
	var options builder.CleanOptions
	for _, o := range opts {
		o(&options)
	}
	// we clean up the bundle by removing its entire conda environment
	cmd := exec.Command("conda", "env", "remove", "-n", b.Source.Name, "-y")
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (p python) String() string {
	return "python"
}

// New returns a python builder that uses conda for dependencies isolation.
func New(opts ...builder.Option) builder.Builder {
	var options builder.Options
	for _, o := range opts {
		o(&options)
	}
	return &python{options: options}
}

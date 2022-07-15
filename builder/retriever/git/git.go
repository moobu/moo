package git

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/moobu/moo/builder/retriever"
)

type git struct {
	options retriever.Options
}

func (g *git) Retrieve(rawURL string, opts ...retriever.RetrieveOption) (*retriever.Repository, error) {
	var options retriever.RetrieveOptions
	for _, o := range opts {
		o(&options)
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	name := u.Path[1:]
	cmd := exec.Command("git", "clone", rawURL, name)
	cmd.Dir = g.options.Dir
	// TODO: open a new session logger for the output
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	path := filepath.Join(g.options.Dir, name)
	ref := options.Ref
	if len(ref) > 0 {
		// The hard man gave me a lot of lessons woo...
		cmd := exec.Command("git", "reset", "--hard", ref)
		cmd.Dir = path
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}

	return &retriever.Repository{Path: path}, nil
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

func rh4() string {
	var buf [4]byte
	rand.Read(buf[:])
	return fmt.Sprintf("%x", buf)
}

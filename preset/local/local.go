package local

import (
	"os"
	"path/filepath"

	"github.com/moobu/moo/builder"
	"github.com/moobu/moo/builder/auto"
	"github.com/moobu/moo/builder/golang"
	"github.com/moobu/moo/builder/python"
	"github.com/moobu/moo/builder/retriever"
	"github.com/moobu/moo/builder/retriever/git"
	"github.com/moobu/moo/internal/cli"
	"github.com/moobu/moo/router"
	"github.com/moobu/moo/router/static"
	"github.com/moobu/moo/runtime"
	"github.com/moobu/moo/runtime/local"
	"github.com/moobu/moo/runtime/local/driver/raw"
	"github.com/moobu/moo/server"
	"github.com/moobu/moo/server/http"
)

type Preset struct{}

func (Preset) Setup(c cli.Ctx) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	rootdir := filepath.Join(home, ".moo")

	runtime.Default = local.New(raw.New())
	router.Default = static.New()
	server.Default = http.New()

	builders := []builder.Builder{
		python.New(),
		golang.New(),
	}
	repodir := filepath.Join(rootdir, "repos")
	builder.Default = retriever.New(
		git.New(retriever.Dir(repodir)),
		auto.New(builders))
	return nil
}

func (Preset) String() string {
	return "local"
}

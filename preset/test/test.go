package test

import (
	"github.com/moobu/moo/builder"
	noopBuilder "github.com/moobu/moo/builder/noop"
	"github.com/moobu/moo/builder/retriever"
	noopRetriever "github.com/moobu/moo/builder/retriever/noop"
	"github.com/moobu/moo/internal/cli"
	"github.com/moobu/moo/router"
	"github.com/moobu/moo/router/static"
	"github.com/moobu/moo/runtime"
	noopRuntime "github.com/moobu/moo/runtime/noop"
	"github.com/moobu/moo/server"
	"github.com/moobu/moo/server/http"
)

type Preset struct{}

func (Preset) Setup(c cli.Ctx) error {
	builder.Default = retriever.New(noopRetriever.New(), noopBuilder.New())
	runtime.Default = noopRuntime.New()
	router.Default = static.New()
	server.Default = http.New()
	return nil
}

func (Preset) String() string {
	return "test"
}

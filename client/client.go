package client

import (
	"github.com/moobu/moo/builder"
	"github.com/moobu/moo/router"
	"github.com/moobu/moo/runtime"
)

type Client interface {
	runtime.Runtime
	builder.Builder
	router.Router
}

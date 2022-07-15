package kubernetes

import (
	"github.com/moobu/moo/internal/cli"
	"github.com/moobu/moo/router"
	"github.com/moobu/moo/runtime"
	"github.com/moobu/moo/runtime/kubernetes"
	"github.com/moobu/moo/server"
	"github.com/moobu/moo/server/http"
)

type Preset struct{}

func (Preset) Setup(c cli.Ctx) error {
	runtime.Default = kubernetes.New()
	router.Default = nil // TODO: should we have a router built on kubernetes?
	server.Default = http.New()
	return nil
}

func (Preset) String() string {
	return "kubernetes"
}

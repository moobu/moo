package noop

import (
	"fmt"
	"net/http"

	"github.com/moobu/moo/gateway"
)

type noop struct {
	options gateway.ProxyOptions
}

func (n noop) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// for now, we only allow POST requests
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// strip out the protocol name
	pod := r.URL.Path[5:]
	fmt.Fprint(w, pod)
}

func (n noop) String() string {
	return "noop"
}

func New(opts ...gateway.ProxyOption) gateway.Proxy {
	var options gateway.ProxyOptions
	for _, o := range opts {
		o(&options)
	}
	return &noop{options: options}
}

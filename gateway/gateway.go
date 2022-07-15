package gateway

import (
	"fmt"
	"net"
	"net/http"
)

// Gateway passes on the incomming requests to the pods using
// registered Proxies. It has only one default implementation.
type Gateway interface {
	Handle(Proxy) error
	Serve(net.Listener) error
}

// Proxy is used by the gateway to pass on the incomming requests
// to the pod using the same protocol the Proxy uses.
type Proxy interface {
	http.Handler
	String() string
}

// here is the only gateway implementation needed
type gateway struct {
	options Options
	mux     *http.ServeMux
}

func (g *gateway) Handle(proxy Proxy) error {
	g.mux.Handle(fmt.Sprintf("/%s/", proxy.String()), Cors(proxy))
	return nil
}

func (g *gateway) Serve(l net.Listener) error {
	return http.Serve(l, g.mux)
}

func New(opts ...Option) Gateway {
	var options Options
	for _, o := range opts {
		o(&options)
	}
	return &gateway{
		options: options,
		mux:     http.NewServeMux(),
	}
}

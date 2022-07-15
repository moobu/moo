package http

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/moobu/moo/server"
)

type httpServer struct {
	options server.Options
}

func (s *httpServer) Serve(l net.Listener) error {
	mux := http.NewServeMux()
	// runtime
	mux.HandleFunc("/create", Create)
	mux.HandleFunc("/delete", Delete)
	mux.HandleFunc("/list", List)
	// routes
	mux.HandleFunc("/register", Register)
	mux.HandleFunc("/deregister", Deregister)
	mux.HandleFunc("/lookup", Lookup)
	// builder
	mux.HandleFunc("/build", Build)
	mux.HandleFunc("/clean", Clean)
	return http.Serve(l, mux)
}

func (httpServer) String() string {
	return "http"
}

func New(opts ...server.Option) server.Server {
	var options server.Options
	for _, o := range opts {
		o(&options)
	}
	return &httpServer{options: options}
}

const headerError = "X-Moo-Error"

func writeJSON(w http.ResponseWriter, v any, err error) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	if err != nil {
		header.Set(headerError, err.Error())
		return
	}
	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

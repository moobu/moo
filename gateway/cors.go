package gateway

import "net/http"

// TODO: implement the cors wrapper
func Cors(h http.Handler) http.Handler {
	return h
}

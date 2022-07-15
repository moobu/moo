package server

import "net"

type Server interface {
	Serve(net.Listener) error
	String() string
}

var Default Server

func Serve(l net.Listener) error {
	return Default.Serve(l)
}

func String() string {
	return Default.String()
}

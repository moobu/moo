package logger

// TODO: design
type Logger interface {
	Open(*Object) (Stream, error)
}

type Stream interface{}

type Object struct {
	Name   string
	Module string
}

var Default Logger

func Open(o *Object) (Stream, error) {
	return Default.Open(o)
}

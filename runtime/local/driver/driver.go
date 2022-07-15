package driver

import (
	"io"

	"github.com/moobu/moo/builder"
)

type Driver interface {
	// Fork creates a new process group
	Fork(*Runnable) (*Process, error)
	// Kill kills the process and its child processes
	Kill(*Process) error
	// Wait waits the process to exit
	Wait(*Process) error
}

type Runnable struct {
	Bundle *builder.Bundle
	Env    []string
	Args   []string
}

type Process struct {
	ID  int
	In  io.Writer
	Out io.Reader
	Err io.Reader
}

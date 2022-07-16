package cluster

import "github.com/moobu/moo/internal/cli"

type Preset struct{}

func (Preset) Setup(c cli.Ctx) error {
	return nil
}

func (Preset) String() string {
	return "cluster"
}

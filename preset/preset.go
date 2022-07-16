package preset

import (
	"errors"

	"github.com/moobu/moo/internal/cli"
	"github.com/moobu/moo/preset/cluster"
	"github.com/moobu/moo/preset/local"
	"github.com/moobu/moo/preset/test"
)

type Preset interface {
	Setup(cli.Ctx) error
	String() string
}

var preset = map[string]Preset{
	"test":    test.Preset{},
	"local":   local.Preset{},
	"cluster": cluster.Preset{},
}

func Register(p Preset) {
	preset[p.String()] = p
}

func Deregister(p Preset) {
	delete(preset, p.String())
}

func Use(c cli.Ctx, name string) error {
	preset, ok := preset[name]
	if !ok {
		return errors.New("no such preset")
	}
	return preset.Setup(c)
}

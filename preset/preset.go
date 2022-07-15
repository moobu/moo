package preset

import (
	"errors"

	"github.com/moobu/moo/internal/cli"
	"github.com/moobu/moo/preset/kubernetes"
	"github.com/moobu/moo/preset/local"
	"github.com/moobu/moo/preset/test"
)

type Presets interface {
	Setup(cli.Ctx) error
	String() string
}

var preset = map[string]Presets{
	"test":       test.Preset{},
	"local":      local.Preset{},
	"kubernetes": kubernetes.Preset{},
}

func Register(p Presets) {
	preset[p.String()] = p
}

func Deregister(p Presets) {
	delete(preset, p.String())
}

func Use(c cli.Ctx, name string) error {
	preset, ok := preset[name]
	if !ok {
		return errors.New("no such preset")
	}
	return preset.Setup(c)
}

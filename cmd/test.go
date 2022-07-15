package cmd

import (
	"fmt"

	"github.com/moobu/moo/internal/cli"
)

func init() {
	cmd.Register(&cli.Cmd{
		Name:  "test",
		About: "test subcommand",
		Pos:   []string{"a", "b"},
		Run: func(c cli.Ctx) error {
			fmt.Println(c.Pos())
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "x",
			},
			&cli.IntFlag{
				Name: "y",
			},
			&cli.BoolFlag{
				Name: "z",
			},
		},
	})
}

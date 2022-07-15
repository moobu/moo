package cmd

import "github.com/moobu/moo/internal/cli"

func init() {
	cmd.Register(&cli.Cmd{
		Name:  "logs",
		About: "output log file",
		Run:   Logs,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "watch",
				Usage: "enable live stream",
			},
			&cli.StringFlag{
				Name:  "ns",
				Usage: "specify a namespace to act on",
				Value: defaultNamespace,
			},
		},
	})
}

func Logs(c cli.Ctx) error {
	return nil
}

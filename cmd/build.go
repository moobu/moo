package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/moobu/moo/builder"
	"github.com/moobu/moo/client"
	"github.com/moobu/moo/client/http"
	"github.com/moobu/moo/internal/cli"
)

func init() {
	cmd.Register(&cli.Cmd{
		Name:  "build",
		About: "build a source",
		Pos:   []string{"source"},
		Run:   Build,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "specify a name for the deploy",
			},
			&cli.StringFlag{
				Name:  "ref",
				Usage: "specify a reference to deploy",
			},
			&cli.StringFlag{
				Name:  "builder",
				Usage: "specify a builder to build the source",
				Value: "python",
			},
			&cli.BoolFlag{
				Name:  "local",
				Usage: "indicate that a local source is being used",
			},
			&cli.StringFlag{
				Name:  "server",
				Usage: "specify the address of Moo server",
				Value: defaultServerAddr,
			},
			&cli.StringFlag{
				Name:  "ns",
				Usage: "specify a namespace to act on",
				Value: defaultNamespace,
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "enable JSON formatted output",
			},
		},
	})
}

func Build(c cli.Ctx) error {
	cli := http.New(client.Server(c.String("server")))
	rawURL := c.Pos()[0]
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	name := c.String("name")
	if len(name) == 0 {
		name = filepath.Base(u.Path)
	}

	ns := c.String("ns")
	source := &builder.Source{
		Name: name,
		URL:  rawURL,
		Type: c.String("builder"),
	}
	bundle, err := build(c, cli, source, ns)
	if err != nil {
		return err
	}

	if c.Bool("json") {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "    ")
		return enc.Encode(bundle)
	}

	fmt.Println("Done.")
	return nil
}

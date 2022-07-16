package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/moobu/moo/builder"
	"github.com/moobu/moo/client"
	"github.com/moobu/moo/client/http"
	"github.com/moobu/moo/internal/cli"
	"github.com/moobu/moo/runtime"
)

func init() {
	cmd.Register(&cli.Cmd{
		Name:  "run",
		About: "run a pod",
		Pos:   []string{"source"},
		Run:   Run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "specify a name for the deploy",
			},
			&cli.StringFlag{
				Name:  "tag",
				Usage: "specify a tag for the deploy",
			},
			&cli.StringFlag{
				Name:  "ref",
				Usage: "specify a commit reference of the source",
			},
			&cli.StringFlag{
				Name:  "builder",
				Usage: "specify a builder to build the source",
				Value: "python",
			},
			&cli.IntFlag{
				Name:  "replicas",
				Usage: "specify replicas to be deployed",
				Value: 1,
			},
			&cli.BoolFlag{
				Name:  "gpu",
				Usage: "enable GPU support for the pod",
			},
			&cli.BoolFlag{
				Name:  "output",
				Usage: "enable the client side stdard output",
				Value: true,
			},
			&cli.StringSliceFlag{
				Name:  "env",
				Usage: "specify env variables for the deploy",
			},
			&cli.StringFlag{
				Name:  "server",
				Usage: "specify the address of Moo server",
				Value: defaultServerAddr,
			},
			&cli.BoolFlag{
				Name:  "image",
				Usage: "indicate that an pre-built image is being used",
			},
			&cli.BoolFlag{
				Name:  "local",
				Usage: "indicate that a local source is being used",
			},
			&cli.StringFlag{
				Name:  "ns",
				Usage: "specify a namespace to act on",
				Value: defaultNamespace,
			},
		},
	})
}

// TODO: how can we watch the log output by the server end?
func Run(c cli.Ctx) error {
	// connect to the Moo server
	cli := http.New(client.Server(c.String("server")))
	// parse the given remote address of the source.
	rawURL := c.Pos()[0]
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	// use the base name of the source if no name is given.
	name := c.String("name")
	if len(name) == 0 {
		name = filepath.Base(u.Path)
	}
	// tag the deploy as 'latest' by default
	tag := c.String("tag")
	if len(tag) == 0 {
		tag = "latest"
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

	// specify the output
	// TODO: use a default file on the machine running the CLI.
	output := io.Discard
	if c.Bool("output") {
		output = os.Stdout
	}
	pod := &runtime.Pod{
		Name: name,
		Tag:  tag,
	}
	// build the optional functions used by the Moo runtime.
	opts := []runtime.CreateOption{
		runtime.CreateNamespace(ns),
		runtime.Env(c.StringSlice("env")...),
		runtime.Replicas(c.Int("replicas")),
		runtime.GPU(c.Bool("gpu")),
		runtime.Output(output),
		runtime.Bundle(bundle),
	}
	// tell the server to run a pod containing the bundle.
	if err := cli.Create(pod, opts...); err != nil {
		return err
	}
	fmt.Printf("Deployed.\n")
	return nil
}

func build(c cli.Ctx, cli client.Client, s *builder.Source, ns string) (*builder.Bundle, error) {
	// we don't need to call the builder if we are using
	// a pre-built image that can be retrieved from remote.
	if c.Bool("image") {
		return &builder.Bundle{Source: s}, nil
	}
	// TODO: handle local sources
	if c.Bool("local") {
		return nil, errors.New("local sources are unsupported yet")
	}
	// tell the server to build the source, so we can then use
	// the runtime to run the returned bundle.
	return cli.Build(s,
		builder.Ref(c.String("ref")),
		builder.BuildNamespace(ns))
}

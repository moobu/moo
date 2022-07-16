package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/moobu/moo/client"
	"github.com/moobu/moo/client/http"
	"github.com/moobu/moo/internal/cli"
	"github.com/moobu/moo/runtime"
)

func init() {
	cmd.Register(&cli.Cmd{
		Name:  "ps",
		About: "list running pods",
		Run:   List,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "all",
				Usage: "list all",
			},
			&cli.StringFlag{
				Name:  "name",
				Usage: "specify the pod's name",
			},
			&cli.StringFlag{
				Name:  "tag",
				Usage: "specify the pod's tag",
			},
			&cli.StringFlag{
				Name:  "server",
				Usage: "specify the address of Moo server",
				Value: defaultServerAddr,
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "enable JSON formatted output",
			},
			&cli.StringFlag{
				Name:  "ns",
				Usage: "specify a namespace to act on",
				Value: defaultNamespace,
			},
		},
	})
}

func List(c cli.Ctx) error {
	cli := http.New(client.Server(c.String("server")))
	pods, err := cli.List(
		runtime.Name(c.String("name")),
		runtime.Tag(c.String("tag")),
		runtime.ListNamespace(c.String("ns")))
	if err != nil {
		return err
	}

	if c.Bool("json") {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "   ")
		return enc.Encode(pods)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprint(tw, "NAME\tTAG\tSTATUS\tUPTIME\tSOURCE\n")

	for _, pod := range pods {
		updated := pod.Get("updated")
		status := pod.Get("status")
		source := pod.Get("source")
		uptime := updated

		if started, err := time.Parse(time.RFC3339, updated); err == nil {
			uptime = time.Since(started).Round(time.Second).String()
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", pod.Name, pod.Tag, status, uptime, source)
	}
	return tw.Flush()
}

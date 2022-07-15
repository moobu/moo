package main

import (
	"context"
	"fmt"
	"os"

	"github.com/moobu/moo/cmd"
)

func main() {
	ctx := context.Background()
	if err := cmd.RunCtx(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

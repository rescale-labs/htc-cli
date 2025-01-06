package main

import (
	"context"
	"github.com/rescale-labs/htc-cli/cli"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		cli.Usage()
	}
	transferOpts, err := cli.ParseArgs(os.Args[2:])
	if err != nil {
		cli.Usage()
	}

	ctx := context.Background()
	err = transferOpts.Transfer(ctx)
	if err != nil {
		log.Fatalf("error running command %v", err)
	}
}

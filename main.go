package main

import (
	"context"
	"github.com/rescale/htc-storage-cli/cli"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		cli.Usage()
	}
	transferOpts, err := cli.ParseArgs(os.Args[2:])
	if err != nil {
		log.Fatalf("error parsing arguments %v", err)
	}

	ctx := context.Background()
	err = transferOpts.Transfer(ctx)
	if err != nil {
		log.Fatalf("error running command %v", err)
	}
}

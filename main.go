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
	src, dest := cli.ParseArgs(os.Args[2:])

	ctx := context.Background()
	transfer := cli.NewTransfer(src, dest)
	err := transfer.Transfer(ctx)
	if err != nil {
		log.Fatalf("error running command %s", err)
	}
}

package main

import (
	"os"

	"github.com/rescale-labs/htc-cli/v2/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		// fmt.Fprintf(os.Stderr, "Command failed: %v", err)
		os.Exit(1)
	}
}

package common

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rescale-labs/htc-cli/v2/config"
)

func WrapRunE(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		err := f(cmd, args)
		if err != nil {
			if _, ok := err.(*config.UsageError); ok {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: %s\n\n", err)
				cmd.Usage()
				os.Exit(1)
			}
		}
		cobra.CheckErr(err)
	}
}

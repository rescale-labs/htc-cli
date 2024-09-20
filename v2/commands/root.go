package commands

import (
	"github.com/spf13/cobra"

	"github.com/rescale/htc-storage-cli/v2/commands/auth"
	"github.com/rescale/htc-storage-cli/v2/commands/project"
)

var RootCmd = &cobra.Command{
	Use:   "htc",
	Short: "The CLI for Rescale's High Throughput Computing (HTC) API",
	Long: `htc provides easy access to key parts of Rescale's HTC API.
  	See https://htc.rescale.com/docs/ for more details.`,
}

func init() {
	// cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringP("output", "o", "yaml", "output format")

	authCmd := &cobra.Command{
		Use: "auth",
	}
	authCmd.AddCommand(auth.LoginCmd)
	authCmd.AddCommand(auth.WhoAmICmd)
	RootCmd.AddCommand(authCmd)

	projectCmd := &cobra.Command{
		Use: "project",
	}
	projectCmd.AddCommand(project.GetCmd)
	projectCmd.AddCommand(project.LimitsCmd)
	RootCmd.AddCommand(projectCmd)
}

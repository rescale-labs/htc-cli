package commands

import (
	"github.com/spf13/cobra"

	"github.com/rescale/htc-storage-cli/v2/commands/auth"
	"github.com/rescale/htc-storage-cli/v2/commands/metrics"
	"github.com/rescale/htc-storage-cli/v2/commands/project"
	"github.com/rescale/htc-storage-cli/v2/commands/task"
)

var RootCmd = &cobra.Command{
	Use:   "htc",
	Short: "The CLI for Rescale's High Throughput Computing (HTC) API",
	Long: `htc provides easy access to key parts of Rescale's HTC API.
  	See https://htc.rescale.com/docs/ for more details.`,
}

func init() {
	// cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringP("output", "o", "", "output format")

	authCmd := &cobra.Command{
		Use: "auth",
	}
	authCmd.AddCommand(auth.LoginCmd)
	authCmd.AddCommand(auth.WhoAmICmd)
	RootCmd.AddCommand(authCmd)

	metricsCmd := &cobra.Command{
		Use: "metrics",
	}
	metricsCmd.AddCommand(metrics.GetCmd)
	RootCmd.AddCommand(metricsCmd)

	projectCmd := &cobra.Command{
		Use: "project",
	}
	projectCmd.AddCommand(project.GetCmd)
	projectCmd.AddCommand(project.LimitsCmd)
	RootCmd.AddCommand(projectCmd)

	taskCmd := &cobra.Command{
		Use: "task",
	}
	taskCmd.AddCommand(task.GetCmd)
	RootCmd.AddCommand(taskCmd)

}

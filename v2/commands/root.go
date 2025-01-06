package commands

import (
	"github.com/spf13/cobra"

	"github.com/rescale-labs/htc-cli/v2/commands/auth"
	"github.com/rescale-labs/htc-cli/v2/commands/config"
	"github.com/rescale-labs/htc-cli/v2/commands/image"
	"github.com/rescale-labs/htc-cli/v2/commands/job"
	"github.com/rescale-labs/htc-cli/v2/commands/metrics"
	"github.com/rescale-labs/htc-cli/v2/commands/project"
	"github.com/rescale-labs/htc-cli/v2/commands/region"
	"github.com/rescale-labs/htc-cli/v2/commands/task"
	"github.com/rescale-labs/htc-cli/v2/commands/version"
)

var RootCmd = &cobra.Command{
	Use:   "htc",
	Short: "The CLI for Rescale's High Throughput Computing (HTC) API",
	Long: `htc provides easy access to key parts of Rescale's HTC API.
See https://htc.rescale.com/docs/ for more details.`,
}

func init() {
	// cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringP("output", "o", "text", "output format")
	RootCmd.PersistentFlags().String("context", "", "config context to use (default \"default\")")

	// auth
	RootCmd.AddCommand(auth.AuthCmd)

	// config
	RootCmd.AddCommand(config.ConfigCmd)

	// image
	RootCmd.AddCommand(image.ImageCmd)

	// job
	RootCmd.AddCommand(job.JobCmd)

	// metrics
	RootCmd.AddCommand(metrics.MetricsCmd)

	// project
	RootCmd.AddCommand(project.ProjectCmd)

	// region
	RootCmd.AddCommand(region.RegionCmd)

	// task
	RootCmd.AddCommand(task.TaskCmd)

	// version
	RootCmd.AddCommand(version.Cmd)
}

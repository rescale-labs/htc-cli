package commands

import (
	"github.com/spf13/cobra"

	"github.com/rescale/htc-storage-cli/v2/commands/auth"
	"github.com/rescale/htc-storage-cli/v2/commands/config"
	cfgcontext "github.com/rescale/htc-storage-cli/v2/commands/config/context"
	"github.com/rescale/htc-storage-cli/v2/commands/image"
	"github.com/rescale/htc-storage-cli/v2/commands/job"
	"github.com/rescale/htc-storage-cli/v2/commands/metrics"
	"github.com/rescale/htc-storage-cli/v2/commands/project"
	"github.com/rescale/htc-storage-cli/v2/commands/region"
	"github.com/rescale/htc-storage-cli/v2/commands/task"
	"github.com/rescale/htc-storage-cli/v2/commands/version"
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
	configCmd := &cobra.Command{
		Use: "config",
	}
	configCmd.AddCommand(config.SetCmd)
	configCmd.AddCommand(config.UnsetCmd)
	RootCmd.AddCommand(configCmd)

	// config context
	contextCmd := &cobra.Command{
		Use: "context",
	}
	contextCmd.AddCommand(cfgcontext.DeleteCmd)
	contextCmd.AddCommand(cfgcontext.GetCmd)
	contextCmd.AddCommand(cfgcontext.UseCmd)
	configCmd.AddCommand(contextCmd)

	// image
	RootCmd.AddCommand(image.ImageCmd)

	// job
	jobCmd := &cobra.Command{
		Use: "job",
	}
	jobCmd.AddCommand(job.SubmitCmd)
	jobCmd.AddCommand(job.GetCmd)
	RootCmd.AddCommand(jobCmd)

	// metrics
	metricsCmd := &cobra.Command{
		Use: "metrics",
	}
	metricsCmd.AddCommand(metrics.GetCmd)
	RootCmd.AddCommand(metricsCmd)

	// project
	RootCmd.AddCommand(project.ProjectCmd)

	// region
	RootCmd.AddCommand(region.RegionCmd)

	// task
	taskCmd := &cobra.Command{
		Use: "task",
	}
	taskCmd.AddCommand(task.CreateCmd)
	taskCmd.AddCommand(task.GetCmd)
	RootCmd.AddCommand(taskCmd)

	// version
	RootCmd.AddCommand(version.Cmd)
}

package commands

import (
	"github.com/spf13/cobra"

	"github.com/rescale/htc-storage-cli/v2/commands/auth"
	"github.com/rescale/htc-storage-cli/v2/commands/image"
	"github.com/rescale/htc-storage-cli/v2/commands/job"
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

	// auth
	authCmd := &cobra.Command{
		Use: "auth",
	}
	authCmd.AddCommand(auth.LoginCmd)
	authCmd.AddCommand(auth.WhoAmICmd)
	RootCmd.AddCommand(authCmd)

	// image
	imageCmd := &cobra.Command{
		Use: "image",
	}
	imageCmd.AddCommand(image.GetCmd)
	// imageCmd.AddCommand(image.CreateRepoCmd)
	// imageCmd.AddCommand(image.LoginCmd)
	RootCmd.AddCommand(imageCmd)

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
	projectCmd := &cobra.Command{
		Use: "project",
	}
	projectCmd.AddCommand(project.DimensionsCmd)
	projectCmd.AddCommand(project.GetCmd)
	projectCmd.AddCommand(project.LimitsCmd)
	RootCmd.AddCommand(projectCmd)

	// task
	taskCmd := &cobra.Command{
		Use: "task",
	}
	taskCmd.AddCommand(task.CreateCmd)
	taskCmd.AddCommand(task.GetCmd)
	RootCmd.AddCommand(taskCmd)

}

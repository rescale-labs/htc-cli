package project

import (
	"github.com/spf13/cobra"
)

var ProjectCmd = &cobra.Command{
	Use: "project",
}

func init() {
	ProjectCmd.PersistentFlags().String("project-id", "", "HTC project ID")

	ProjectCmd.AddCommand(DimensionsCmd)
	ProjectCmd.AddCommand(GetCmd)
	ProjectCmd.AddCommand(LimitsCmd)
	ProjectCmd.AddCommand(RetentionPolicyCmd)
}

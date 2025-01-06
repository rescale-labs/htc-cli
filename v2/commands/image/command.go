package image

import (
	"github.com/spf13/cobra"
)

var ImageCmd = &cobra.Command{
	Use:   "image",
	Short: "Commands for managing container images",
}

func init() {
	ImageCmd.PersistentFlags().String("project-id", "", "HTC project ID")
	ImageCmd.AddCommand(GetCmd)
	ImageCmd.AddCommand(CreateRepoCmd)
	ImageCmd.AddCommand(LoginRepoCmd)
	ImageCmd.AddCommand(PushCmd)
}

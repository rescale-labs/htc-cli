package image

import (
	"github.com/spf13/cobra"
)

var ImageCmd = &cobra.Command{
	Use: "image",
}

func init() {
	ImageCmd.PersistentFlags().String("project-id", "", "HTC project ID")
	ImageCmd.AddCommand(GetCmd)
	ImageCmd.AddCommand(CreateRepoCmd)
	ImageCmd.AddCommand(LoginRepoCmd)
}

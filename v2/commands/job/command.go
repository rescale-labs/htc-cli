package job

import (
	"github.com/spf13/cobra"
)

var JobCmd = &cobra.Command{
	Use:   "job",
	Short: "Commands for HTC projects",
}

func init() {
	JobCmd.AddCommand(SubmitCmd)
	JobCmd.AddCommand(GetCmd)
}

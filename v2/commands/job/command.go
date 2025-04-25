package job

import (
	"github.com/spf13/cobra"
)

var JobCmd = &cobra.Command{
	Use:   "job",
	Short: "Commands for HTC jobs",
}

func init() {
	JobCmd.AddCommand(SubmitCmd)
	JobCmd.AddCommand(GetCmd)
	JobCmd.AddCommand(CancelCmd)
	JobCmd.AddCommand(LogsCmd)
	JobCmd.AddCommand(EventsCmd)
}

package task

import (
	"github.com/spf13/cobra"
)

var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "Commands for managing HTC tasks",
}

func init() {
	TaskCmd.AddCommand(CreateCmd)
	TaskCmd.AddCommand(GetCmd)
	TaskCmd.AddCommand(StatsCmd)
}

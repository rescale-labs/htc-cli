package config

import (
	cfgcontext "github.com/rescale/htc-storage-cli/v2/commands/config/context"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage local configuration",
}

func init() {
	ConfigCmd.AddCommand(SetCmd)
	ConfigCmd.AddCommand(UnsetCmd)

	contextCmd := &cobra.Command{
		Use: "context",
	}
	contextCmd.AddCommand(cfgcontext.DeleteCmd)
	contextCmd.AddCommand(cfgcontext.GetCmd)
	contextCmd.AddCommand(cfgcontext.UseCmd)
	ConfigCmd.AddCommand(contextCmd)
}

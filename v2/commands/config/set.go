package config

import (
	"github.com/spf13/cobra"

	"github.com/rescale/htc-storage-cli/v2/common"
	"github.com/rescale/htc-storage-cli/v2/config"
)

func set(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}

	global, err := cmd.Flags().GetBool("global")
	if err != nil {
		return config.UsageErrorf("Error parsing --global: %w", err)
	}

	if err := runner.Config.Set(args[0], args[1], global); err != nil {
		return err
	}
	return nil
}

func unset(cmd *cobra.Command, args []string) error {
	return set(cmd, []string{args[0], ""})
}

var SetCmd = &cobra.Command{
	Use:   "set KEY VALUE",
	Short: "Sets a configuration key to a given value.",
	Long: `Available keys:
  api_url: Rescale HTC API URL
  project_id: Default HTC project ID (can't be set globally)
  task_id: Default HTC task ID (can't be set globally)
`,
	Args: cobra.ExactArgs(2),
	Run:  common.WrapRunE(set),
}

var UnsetCmd = &cobra.Command{
	Use:   "unset KEY",
	Short: "Unsets a configuration key to a given value.",
	Long: `Available keys:
  api_url: Rescale HTC API URL
  project_id: Default HTC project ID (can't be set globally)
  task_id: Default HTC task ID (can't be set globally)
`,
	Args: cobra.ExactArgs(1),
	Run:  common.WrapRunE(unset),
}

func init() {
	SetCmd.Flags().Bool("global", false, "Apply to global config section")
	UnsetCmd.Flags().Bool("global", false, "Apply to global config section")
}

package context

import (
	"github.com/spf13/cobra"

	"github.com/rescale-labs/htc-cli/v2/common"
	"github.com/rescale-labs/htc-cli/v2/config"
)

func deleteCmd(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}
	contextName := args[0]
	if runner.Config.Context == contextName {
		runner.Config.Set("selected_context", config.DefaultContextName, true)
	}
	return runner.Config.Delete(args[0])
}

var DeleteCmd = &cobra.Command{
	Use:   "delete CONTEXT_NAME",
	Short: "Deletes the given config context",
	Args:  cobra.ExactArgs(1),
	Run:   common.WrapRunE(deleteCmd),
}

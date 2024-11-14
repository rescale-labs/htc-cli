package context

import (
	"github.com/spf13/cobra"

	"github.com/rescale/htc-storage-cli/v2/common"
)

func use(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}

	return runner.Config.Set("selected_context", args[0], true)
}

var UseCmd = &cobra.Command{
	Use:   "use CONTEXT_NAME",
	Short: "Sets the current config context",
	Long: `Sets the specified context as the current context. 
        This allows you to run commands without specifying the --project-id and --task-id flags.`,
	Args:  cobra.ExactArgs(1),
	Run:   common.WrapRunE(use),
}

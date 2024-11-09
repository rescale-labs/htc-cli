package context

import (
	"cmp"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/rescale/htc-storage-cli/v2/common"
	"github.com/rescale/htc-storage-cli/v2/tabler"
)

func Get(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}

	g, err := runner.Config.ReadGlobalConf()
	if err != nil {
		return err
	}

	var items tabler.ContextConfs
	for k, v := range g.Contexts {
		items = append(items, &tabler.ContextConf{
			Name:        k,
			Selected:    k == runner.Config.Context,
			ContextConf: v,
		})
	}

	slices.SortFunc(items, func(l, r *tabler.ContextConf) int {
		return cmp.Compare(l.Name, r.Name)
	})

	return runner.PrintResult(items, os.Stdout)
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns all config contexts. Currently selected one is marked by`*`.",
	// Long:
	Run: common.WrapRunE(Get),
}

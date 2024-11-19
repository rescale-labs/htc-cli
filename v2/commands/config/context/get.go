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
		item := &tabler.ContextConf{
			Name:        k,
			Selected:    k == runner.Config.Context,
			ContextConf: v,
		}
		if err := runner.Config.ReadIdentity(k, &item.Identity); err != nil {
			return err
		}
		items = append(items, item)
	}

	slices.SortFunc(items, func(l, r *tabler.ContextConf) int {
		return cmp.Compare(l.Name, r.Name)
	})

	return runner.PrintResult(items, os.Stdout)
}

var GetCmd = &cobra.Command{
	Use: "get",
	Short: `Returns all config contexts. Currently selected one is marked by '*'.

To view all data, use '-o json' or '-o yaml'`,
	Run: common.WrapRunE(Get),
}

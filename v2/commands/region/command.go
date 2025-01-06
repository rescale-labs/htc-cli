package region

import (
	"os"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
	"github.com/rescale-labs/htc-cli/v2/tabler"
	"github.com/spf13/cobra"
)

func getRegions(cmd *cobra.Command, _ []string) error {
	runner, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}

	var r oapi.RescaleRegion
	var res []string
	for _, s := range r.AllValues() {
		res = append(res, string(s))
	}
	return runner.PrintResult(tabler.Regions(res), os.Stdout)
}

var RegionCmd = &cobra.Command{
	Use:   "region",
	Short: "List HTC regions",
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Lists HTC all regions.",
	Long:  "Lists HTC all regions. Not all regions may be available to this workspace.",
	Run:   common.WrapRunE(getRegions),
	Args:  cobra.ExactArgs(0),
}

func init() {
	RegionCmd.AddCommand(GetCmd)
}

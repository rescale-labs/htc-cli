package project

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/common"
)

func LimitsGet(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.GetLimits(ctx,
		oapi.GetLimitsParams{ProjectId: p.ProjectId})
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.HTCProjectLimits:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.GetLimitsForbidden, *oapi.GetLimitsUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}

	return fmt.Errorf("Unknown response type: %s", res)
}

var LimitsCmd = &cobra.Command{
	Use:   "limits",
	Short: "Commands for project limits",
}

var LimitsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns limits for an HTC project",
	Run:   common.WrapRunE(LimitsGet),
	Args:  cobra.ExactArgs(0),
}

// var LimitsApplyCmd = &cobra.Command{
// 	Use:   "apply",
// 	Short: "Sets limits for an HTC project",
// 	Run:   common.WrapRunE(LimitsGet),
// 	Args:  cobra.ExactArgs(0),
// }
//
// var LimitsDeleteCmd = &cobra.Command{
// 	Use:   "delete",
// 	Short: "Deletes limits for an HTC project",
// 	Run:   common.WrapRunE(LimitsGet),
// 	Args:  cobra.ExactArgs(0),
// }

func init() {
	LimitsCmd.AddCommand(LimitsGetCmd)
}

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

func DimensionsGet(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.GetDimensions(ctx,
		oapi.GetDimensionsParams{ProjectId: p.ProjectId})
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.HTCProjectDimensions:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.GetDimensionsForbidden,
		*oapi.GetDimensionsUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}

	return fmt.Errorf("Unknown response type: %s", res)
}

var DimensionsCmd = &cobra.Command{
	Use:   "dimensions",
	Short: "Commands for project dimensions",
}

var DimensionsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns dimensions for an HTC project",
	Run:   common.WrapRunE(DimensionsGet),
	Args:  cobra.ExactArgs(0),
}

// var DimensionsApplyCmd = &cobra.Command{
// 	Use:   "apply",
// 	Short: "Sets dimensions for an HTC project",
// 	Run:   common.WrapRunE(DimensionsGet),
// 	Args:  cobra.ExactArgs(0),
// }
//
// var DimensionsDeleteCmd = &cobra.Command{
// 	Use:   "delete",
// 	Short: "Deletes dimensions for an HTC project",
// 	Run:   common.WrapRunE(DimensionsGet),
// 	Args:  cobra.ExactArgs(0),
// }

func init() {
	DimensionsCmd.AddCommand(DimensionsGetCmd)
}

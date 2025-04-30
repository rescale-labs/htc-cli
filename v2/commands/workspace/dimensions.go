package workspace

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
	"github.com/rescale-labs/htc-cli/v2/tabler"
)

func DimensionsGet(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireWorkspaceId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.GetWorkspaceDimensions(ctx,
		oapi.GetWorkspaceDimensionsParams{WorkspaceId: p.WorkspaceId})
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.HTCWorkspaceDimensions:
		return runner.PrintResult((*tabler.ComputeEnvs)(res), os.Stdout)
	case *oapi.GetWorkspaceDimensionsForbidden,
		*oapi.GetWorkspaceDimensionsUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}
	return fmt.Errorf("Unknown response type: %s", res)
}

var DimensionsCmd = &cobra.Command{
	Use:   "dimensions",
	Short: "Commands for workspace dimensions",
}

var DimensionsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns dimensions for an HTC workspace",
	Run:   common.WrapRunE(DimensionsGet),
	Args:  cobra.ExactArgs(0),
}

func init() {
	DimensionsCmd.AddCommand(DimensionsGetCmd)
}

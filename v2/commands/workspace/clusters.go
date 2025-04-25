package workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
)

func getClusters(ctx context.Context, c oapi.WorkspaceInvoker, workspaceId string) (*oapi.HTCClusterStatusResponse, error) {
	res, err := c.GetGCPClusters(ctx, oapi.GetGCPClustersParams{
		WorkspaceId: workspaceId,
	})
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *oapi.HTCClusterStatusResponse:
		return res, nil
	case *oapi.GetGCPClustersForbidden,
		*oapi.GetGCPClustersUnauthorized:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Clusters(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	worskpaceId, err := runner.GetWorkspaceId()
	if err != nil {
		return err
	}

	ctx := context.Background()
	res, err := getClusters(ctx, runner.Client, worskpaceId)
	if err != nil {
		return err
	}
	fmt.Print(res)
	return nil
}

var ClustersCmd = &cobra.Command{
	Use: "clusters",
	Short: "Get information about cluster status.",
	Run: common.WrapRunE(Clusters),
}

func init() {
	ClustersCmd.Flags()
}

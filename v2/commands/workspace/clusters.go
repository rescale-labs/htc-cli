package workspace

import (
	"context"
	"fmt"
	"io"
	"os"
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

	p := common.IDParams{RequireWorkspaceId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := getClusters(ctx, runner.Client, p.WorkspaceId)
	if err != nil {
		return err
	}
	return writeRows(res, os.Stdout)
}

func writeRows(clusters *oapi.HTCClusterStatusResponse, w io.Writer) error {
	projectId := clusters.GcpProjectId.Value
	for _, cluster := range clusters.Clusters {
		if _, err := fmt.Fprintf(w, "%-24s %-24s %-24s %-24s %-15s %-24s %-48s \n",
			"PROJECT ID", "CLUSTER NAME", "REGION", "VERSION", "STATUS", "AUTOSCALING", "SUBNETWORK"); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "%-24s %-24s %-24s %-24s %-15s %-24s %-48s\n",
			projectId,
			cluster.Name.Value,
			cluster.Region.Value,
			cluster.Version.Value,
			cluster.Status.Value,
			cluster.Autoscaling.Value,
			cluster.Subnetwork.Value,
		); err != nil {
			return err
		}
		
		if len(cluster.NodePools) > 0 {
			if _, err := fmt.Fprintf(w, "\tNODEPOOLS\n\t %-38s %-38s %-24s %-24s %-24s %15s %15s\n",
				"NAME", "VERSION", "INSTANCE TYPE", "STATUS", "AUTOSCALING", "MIN NODES", "MAX NODES"); err != nil {
				return err
			}
		}

		for _, pool := range cluster.GetNodePools() {
			if _, err := fmt.Fprintf(w, "\t %-38s %-38s %-24s %-24s %-24t %15d %15d\n",
				pool.Name.Value,
				pool.Version.Value,
				pool.InstanceType.Value,
				pool.Status.Value,
				pool.Autoscaling.Value.Enabled.Value,
				pool.Autoscaling.Value.MinNodeCount.Value,
				pool.Autoscaling.Value.MaxNodeCount.Value,
			); err != nil {
				return err
			}
		}
	}
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

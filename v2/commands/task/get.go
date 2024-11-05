package task

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/common"
	"github.com/rescale/htc-storage-cli/v2/config"
	"github.com/rescale/htc-storage-cli/v2/tabler"
)

const pageSize = 500

func getTasks(ctx context.Context, c *oapi.Client, projectId string, pageIndex string) (*oapi.HTCTasksResponse, error) {
	log.Printf("HtcProjectsProjectIdTasksGet: projectId=%s pageIndex=%s pageSize=%d", projectId, pageIndex, pageSize)
	res, err := c.HtcProjectsProjectIdTasksGet(ctx, oapi.HtcProjectsProjectIdTasksGetParams{
		ProjectId: projectId,
		PageIndex: oapi.NewOptString(pageIndex),
		PageSize:  oapi.NewOptInt32(pageSize),
	})
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *oapi.HTCTasksResponse:
		return res, nil
	case *oapi.HtcProjectsProjectIdTasksGetForbidden,
		*oapi.HtcProjectsProjectIdTasksGetUnauthorized:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Get(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return config.UsageErrorf("Error setting limit: %w", err)
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	var items []oapi.HTCTask
	var pageIndex string
	for {
		res, err := getTasks(ctx, runner.Client, p.ProjectId, pageIndex)
		if err != nil {
			return err
		}
		items = append(items, res.Items...)
		if limit > 0 && len(items) >= limit {
			items = items[:limit]
			break
		}

		pageIndex = res.Next.Value.Query().Get("pageIndex")
		if pageIndex == "" {
			break
		}
	}
	return runner.PrintResult(tabler.HTCTasks(items), os.Stdout)
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns HTC tasks in a given project.",
	Run:   common.WrapRunE(Get),
}

func init() {
	GetCmd.Flags().IntP("limit", "l", 0, "Limit response to N items")
	GetCmd.Flags().String("project-id", "", "HTC project ID")
}
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
)

const pageSize = 500

func getTasks(ctx context.Context, c *oapi.Client, projectId string, pageIndex string) (*oapi.HTCTasksResponse, error) {
	log.Printf("HtcProjectsProjectIdTasksGet: pageIndex=%s pageSize=%d", pageIndex, pageSize)
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
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return config.UsageErrorf("Error setting limit: %w", err)
	}

	projectId, err := cmd.Flags().GetString("project-id")
	if err != nil {
		return config.UsageErrorf("Error setting project ID: %w", err)
	}

	runner, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}
	if err := runner.UpdateToken(time.Now()); err != nil {
		return err
	}

	ctx := context.Background()
	var items []oapi.HTCTask
	var pageIndex string
	for {
		res, err := getTasks(ctx, runner.Client, projectId, pageIndex)
		if err != nil {
			return err
		}
		items = append(items, res.Items...)
		if len(items) > limit {
			items = items[:limit]
			break
		}

		pageIndex = res.Next.Value.Query().Get("pageIndex")
		if pageIndex == "" {
			break
		}
	}
	return runner.PrintResult(items, os.Stdout)
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns HTC projects in a workspace",
	Run:   common.WrapRunE(Get),
}

func init() {
	GetCmd.Flags().IntP("limit", "l", 0, "Limit response to N items")
	GetCmd.Flags().String("project-id", "", "HTC project ID")
}

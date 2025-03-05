package project

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
	"github.com/rescale-labs/htc-cli/v2/tabler"
)

func getProjects(ctx context.Context, c oapi.ProjectInvoker, pageIndex string) (*oapi.HTCProjectsResponse, error) {
	log.Printf("HtcProjectsGet: pageIndex=%s pageSize=%d", pageIndex, common.PageSize)
	res, err := c.GetProjects(ctx, oapi.GetProjectsParams{
		OnlyMyProjects: oapi.NewOptBool(false),
		PageIndex:      oapi.NewOptString(pageIndex),
		PageSize:       oapi.NewOptInt32(common.PageSize),
	})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HTCProjectsResponse:
		return res, nil
		// runner.PrintResult(res.Items, os.Stdout)
	case *oapi.GetProjectsForbidden,
		*oapi.GetProjectsUnauthorized:
		return nil, fmt.Errorf("forbidden: %s", res)
	}

	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Get(cmd *cobra.Command, args []string) error {
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return fmt.Errorf("Error setting limit: %w", err)
	}

	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	ctx := context.Background()
	var items []oapi.HTCProject
	var pageIndex string
	for {
		res, err := getProjects(ctx, runner.Client, pageIndex)
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
	return runner.PrintResult(tabler.HTCProjects(items), os.Stdout)
}

var GetCmd = &cobra.Command{
	Use:   "get [PROJECT_ID]",
	Short: "Returns all HTC projects, or a single project, in the current workspace",
	// Long:
	Run: common.WrapRunE(Get),
}

func init() {
	GetCmd.Flags().IntP("limit", "l", 0, "Limit response to N items")
}

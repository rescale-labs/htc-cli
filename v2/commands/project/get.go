package project

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/common"
)

const pageSize = 500

func getProjects(ctx context.Context, c *oapi.Client, pageIndex string) (*oapi.HTCProjectsResponse, error) {
	log.Printf("HtcProjectsGet: pageIndex=%s pageSize=%d", pageIndex, pageSize)
	res, err := c.HtcProjectsGet(ctx, oapi.HtcProjectsGetParams{
		OnlyMyProjects: oapi.NewOptBool(false),
		PageIndex:      oapi.NewOptString(pageIndex),
		PageSize:       oapi.NewOptInt32(pageSize),
	})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HTCProjectsResponse:
		return res, nil
		// runner.PrintResult(res.Items, os.Stdout)
	case *oapi.HtcProjectsGetForbidden,
		*oapi.HtcProjectsGetUnauthorized:
		return nil, fmt.Errorf("forbidden: %s", res)
	}

	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Get(cmd *cobra.Command, args []string) error {
	limit, err := cmd.Flags().GetInt("limit")
	if err != nil {
		return fmt.Errorf("Error setting limit: %w", err)
	}

	runner, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}
	if err := runner.UpdateToken(time.Now()); err != nil {
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
	// Long:
	Run: common.WrapRunE(Get),
}

func init() {
	GetCmd.Flags().IntP("limit", "l", 0, "Limit response to N items")
}

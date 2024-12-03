package job

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

func getJobs(ctx context.Context, c oapi.JobInvoker, params *oapi.GetJobsParams) (*oapi.HTCJobs, error) {
	res, err := c.GetJobs(ctx, *params)
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *oapi.HTCJobs:
		return res, nil
	case *oapi.GetJobsUnauthorized, *oapi.GetJobsForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func getJob(ctx context.Context, c oapi.JobInvoker, projectId, taskId, jobId string) (*oapi.HTCJob, error) {
	params := oapi.GetJobParams{
		ProjectId: projectId,
		TaskId:    taskId,
		JobId:     jobId,
	}
	log.Printf("getJob: %s %s %s", projectId, taskId, jobId)
	res, err := c.GetJob(ctx, params)
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *oapi.HTCJob:
		return res, nil
	case *oapi.GetJobUnauthorized, *oapi.GetJobForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Get(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true, RequireTaskId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	if len(args) == 0 {
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			return config.UsageErrorf("Error setting limit: %w", err)
		}

		group, err := cmd.Flags().GetString("group")
		if err != nil {
			return config.UsageErrorf("Error setting group: %w", err)
		}

		var items []oapi.HTCJob
		params := oapi.GetJobsParams{
			ProjectId: p.ProjectId,
			TaskId:    p.TaskId,
			Group:     oapi.OptString{group, group != ""},
			PageSize:  oapi.NewOptInt32(pageSize),
			ViewType:  oapi.NewOptViewType(oapi.ViewTypeFULL),
		}
		for {
			res, err := getJobs(ctx, runner.Client, &params)
			if err != nil {
				return err
			}
			items = append(items, res.Items...)
			if limit > 0 && len(items) >= limit {
				items = items[:limit]
				break
			}

			params.PageIndex = oapi.NewOptString(
				res.Next.Value.Query().Get("pageIndex"))
			if params.PageIndex.Value == "" {
				break
			}
		}
		return runner.PrintResult(tabler.HTCJobs(items), os.Stdout)
	}

	job, err := getJob(ctx, runner.Client, p.ProjectId, p.TaskId, args[0])
	if err != nil {
		return err
	}
	return runner.PrintResult((*tabler.HTCJob)(job), os.Stdout)
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns HTC jobs in a given task.",
	Run:   common.WrapRunE(Get),
	Args:  cobra.RangeArgs(0, 1),
}

func init() {
	GetCmd.Flags().IntP("limit", "l", 0, "Limit response to N items")
	GetCmd.Flags().String("project-id", "", "HTC project ID")
	GetCmd.Flags().String("task-id", "", "HTC task ID")
	GetCmd.Flags().String("group", "", "HTC job batch group")
}

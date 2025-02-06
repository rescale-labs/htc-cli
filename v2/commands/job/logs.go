package job

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
	"github.com/rescale-labs/htc-cli/v2/config"
	"github.com/rescale-labs/htc-cli/v2/tabler"
)

func logs(ctx context.Context, c oapi.JobInvoker, projectId, taskId, jobId, pageIndex string) (*oapi.HTCJobLogs, error) {
	res, err := c.GetLogs(ctx, oapi.GetLogsParams{
		ProjectId: projectId,
		TaskId:    taskId,
		JobId:     jobId,
		PageSize:  oapi.NewOptInt32(pageSize),
		PageIndex: oapi.NewOptString(pageIndex),
	})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HTCJobLogs:
		return res, nil
	case *oapi.GetLogsUnauthorized,
		*oapi.GetLogsForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Logs(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true, RequireTaskId: true, RequireJobId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	flags := cmd.Flags()
	limit, err := flags.GetInt("limit")
	if err != nil {
		return config.UsageErrorf("Error setting limit: %w", err)
	}

	ctx := context.Background()
	var pageIndex string
	var items []oapi.HTCLogEvent

	for {
		res, err := logs(ctx, runner.Client, p.ProjectId, p.TaskId, p.JobId, pageIndex)
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
	return runner.PrintResult(tabler.HTCJobLog(items), os.Stdout)
}

var LogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Returns HTC job logs given a job ID.",
	Run:   common.WrapRunE(Logs),
}

func init() {
	flags := LogsCmd.Flags()

	flags.IntP("limit", "l", 0, "Limit response to N items")
	flags.String("project-id", "", "HTC project ID")
	flags.String("task-id", "", "HTC task ID")
	flags.String("job-id", "", "HTC job ID")
}

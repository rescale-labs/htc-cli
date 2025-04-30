package task

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

func getTaskStats(ctx context.Context, c oapi.TaskInvoker, projectId, taskId string) (*oapi.JobStatusSummary, error) {
	res, err := c.GetTaskStats(ctx, oapi.GetTaskStatsParams{
		ProjectId: projectId,
		TaskId:    taskId,
	})
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *oapi.JobStatusSummary:
		return res, nil
	case *oapi.GetTaskStatsForbidden,
		*oapi.GetTaskStatsUnauthorized:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}


func Stats(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true, RequireTaskId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}
	
	ctx:= context.Background()

	res, err := getTaskStats(ctx, runner.Client, p.ProjectId, p.TaskId)
	if err != nil {
		return err
	}
	return writeSummary(*res, os.Stdout)
}

func writeSummary(summary oapi.JobStatusSummary, w io.Writer) error {
	if _, err := fmt.Fprintf(w, "%-21s %15s\n", "JOB STATUS", "TOTAL COUNT"); err != nil {
		return err
	}
	
	for status, count := range summary.JobStatuses.Value {
		if _, err := fmt.Fprintf(w, "%-21s %15d\n", status, count); err != nil {
			return err
		}
	}

	return nil
}

var StatsCmd = &cobra.Command{
	Use:   	"stats",
	Short: 	"Get task statistics",
	Run: 	common.WrapRunE(Stats),
}

func init() {
	flags := StatsCmd.Flags()
	flags.String("project-id", "", "HTC project ID")
	flags.String("task-id", "", "HTC task ID")
}

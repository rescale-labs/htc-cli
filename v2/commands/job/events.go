package job

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
	"github.com/rescale-labs/htc-cli/v2/config"
)

func events(ctx context.Context, c oapi.JobInvoker, projectId, taskId, jobId, pageIndex string) (*oapi.HTCJobStatusEvents, error) {
	res, err := c.GetEvents(ctx, oapi.GetEventsParams{
		ProjectId: projectId,
		TaskId:    taskId,
		JobId:     jobId,
		PageSize:  oapi.NewOptInt32(common.PageSize),
		PageIndex: oapi.OptString{pageIndex, pageIndex != ""},
	})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HTCJobStatusEvents:
		return res, nil
	case *oapi.GetEventsUnauthorized,
		*oapi.GetEventsForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %T", res)
}

func Events(cmd *cobra.Command, args[]string) error{
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true, RequireTaskId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	if len(args) != 1 {
		return config.UsageErrorf("Error: job ID not provided")
	}

	jobId := args[0]

	ctx := context.Background()

	// var items []oapi.RescaleJobStatusEvent
	pageIndex := ""
	for {
		res, err := events(ctx, runner.Client, p.ProjectId, p.TaskId, jobId, "")
		if err != nil {
			return err
		}
		// items = append(items, res.Items...)
		fmt.Println(res.Items)
		if len(res.Items) == 0 {
			break
		}

		pageIndex = res.Next.Value.Query().Get("pageIndex")
		if pageIndex == "" {
			break
		}
	}
	// return runner.PrintResult(items, os.Stdout)
	return nil
}

var EventsCmd = &cobra.Command{
	Use: 	"events [JOB_UUID]",
	Short: 	"Returns latest HTC job events given a job ID.",
	Run: 	common.WrapRunE(Events),
	Args: 	cobra.ExactArgs(1),
}

func init() {
	flags := EventsCmd.Flags()
	flags.String("project-id", "", "HTC project ID")
	flags.String("task-id", "", "HTC task ID")
}
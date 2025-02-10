package job

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
)

func cancelJobs(ctx context.Context, c oapi.JobInvoker, params *oapi.CancelJobsParams) (*oapi.CancelJobsOK, error) {
	res, err := c.CancelJobs(ctx, *params)
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *oapi.CancelJobsOK:
		return res, nil
	case *oapi.CancelJobsForbidden:
		return nil, errors.New("make sure you are accessing your own project and task")
	case *oapi.CancelJobsUnauthorized:
		return nil, errors.New("refresh your auth with `htc auth login`")
	default:
		return nil, fmt.Errorf("unknown operation %s", res)
	}
}

func Cancel(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true, RequireTaskId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	group, err := cmd.Flags().GetString("group")
	if err != nil {
		return fmt.Errorf("error setting group: %w", err)
	}

	params := oapi.CancelJobsParams{
		ProjectId: p.ProjectId,
		TaskId:    p.TaskId,
		Group:     oapi.OptString{Value: group, Set: group != ""},
	}

	ctx := context.Background()
	_, err = cancelJobs(ctx, runner.Client, &params)
	if err != nil {
		return err
	}
	return runner.PrintResult("Cancel request sent successfully!", os.Stdout)
}

var CancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancels all jobs in HTC tasks in a given project.",
	Long:  "Cancels all jobs in HTC tasks in a given project. Even when OK is returned cancel is best effort meaning not all jobs are cancelled.",
	Run:   common.WrapRunE(Cancel),
}

func init() {
	flags := CancelCmd.Flags()
	flags.String("project-id", "", "HTC project ID")
	flags.String("task-id", "", "HTC task ID")
	flags.String("group", "", "HTC job batch group")
}

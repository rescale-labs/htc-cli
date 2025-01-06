package task

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
)

func createTask(ctx context.Context, c oapi.TaskInvoker, projectId, name, desc string) (*oapi.HTCTask, error) {
	res, err := c.CreateTask(ctx,
		oapi.NewOptHTCTask(oapi.HTCTask{TaskName: name, TaskDescription: desc}),
		oapi.CreateTaskParams{ProjectId: projectId})
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *oapi.HTCTask:
		return res, nil
	case *oapi.CreateTaskForbidden,
		*oapi.CreateTaskUnauthorized:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Create(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	name := args[0]
	var desc string
	if len(args) == 2 {
		desc = args[1]
	}

	ctx := context.Background()
	res, err := createTask(ctx, runner.Client, p.ProjectId, name, desc)
	if err != nil {
		return err
	}
	return runner.PrintResult(res, os.Stdout)
}

var CreateCmd = &cobra.Command{
	Use:   "create TASK_NAME [TASK_DESCRIPTION]",
	Short: "Creates HTC task in a given project.",
	Args:  cobra.RangeArgs(1, 2),
	Run:   common.WrapRunE(Create),
}

func init() {
	CreateCmd.Flags().String("project-id", "", "HTC project ID")
}

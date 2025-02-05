package job

import (
	"context"
	"fmt"
	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/rescale-labs/htc-cli/v2/common"
)

type cancelRequest struct {
	params oapi.CancelJobsParams
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
		return fmt.Errorf("Error setting group: %w", err)
	}

	ctx := context.Background()
	res, err := runner.Client.CancelJobs(ctx, oapi.CancelJobsParams{p.ProjectId, p.TaskId, oapi.NewOptString(group)})
	if err != nil {
		return err
	}

	switch res.(type) {
	case *oapi.CancelJobsOK:
		return runner.PrintResult("Cancel request sent successfully!", os.Stdout)
	case *oapi.CancelJobsForbidden:
		return runner.PrintResult("Make sure you are accessing your own project and task!", os.Stderr)
	case *oapi.CancelJobsUnauthorized:
		return runner.PrintResult("Refresh your auth!", os.Stderr)
	default:
		return runner.PrintResult(fmt.Errorf("unknown operation %s", res), os.Stderr)
	}
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

package job

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/common"
)

type submitRequest struct {
	batch  []oapi.HTCJobSubmitRequest
	params oapi.HtcProjectsProjectIdTasksTaskIdJobsBatchPostParams
}

func submit(ctx context.Context, c *oapi.Client,
	r *submitRequest) (*oapi.HtcProjectsProjectIdTasksTaskIdJobsBatchPostRes, error) {

	res, err := c.HtcProjectsProjectIdTasksTaskIdJobsBatchPost(ctx, r.batch, r.params)
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	// case *oapi.HTCProjectsResponse:
	// 	return res, nil
	// 	// runner.PrintResult(res.Items, os.Stdout)
	// case *oapi.HtcProjectsGetForbidden,
	// 	*oapi.HtcProjectsGetUnauthorized:
	// 	return nil, fmt.Errorf("forbidden: %s", res)
	case *oapi.HTCRequestError:
		return nil, fmt.Errorf("%s: %s", res.ErrorDescription.Value, res.Message.Value)
	}

	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Submit(cmd *cobra.Command, args []string) error {
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

	var r *os.File
	if len(args) != 1 {
		return fmt.Errorf("Error: job yaml not provided")
	}
	if args[0] == "-" {
		r = os.Stdin
	} else {
		var err error
		r, err = os.Open(args[0])
		if err != nil {
			return fmt.Errorf("Error opening %s: %w", err)
		}
		defer r.Close()
	}

	req := submitRequest{
		params: oapi.HtcProjectsProjectIdTasksTaskIdJobsBatchPostParams{
			ProjectId: p.ProjectId,
			TaskId:    p.TaskId,
			Group:     oapi.NewOptString(group),
		},
	}

	dec := json.NewDecoder(r)
	if err := dec.Decode(&req.batch); err != nil {
		return fmt.Errorf("Error parsing %s: %v", args[0], err)
	}

	ctx := context.Background()
	res, err := submit(ctx, runner.Client, &req)
	if err != nil {
		return fmt.Errorf("Error on job submission: %v", err)
	}
	log.Printf("res: %#v", res)
	return nil
}

var SubmitCmd = &cobra.Command{
	Use:   "submit JOB_BATCH_JSON",
	Args:  cobra.ExactArgs(1),
	Short: "Submits jobs for a given task and project",
	// Long:
	Run: common.WrapRunE(Submit),
}

func init() {
	SubmitCmd.Flags().String("project-id", "", "HTC project ID (required)")
	SubmitCmd.Flags().String("task-id", "", "HTC task ID (required)")
	SubmitCmd.Flags().String("group", "", "Group")
}

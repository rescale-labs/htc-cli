package job

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-faster/yaml"
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

	// switch res := res.(type) {
	// case *oapi.HTCProjectsResponse:
	// 	return res, nil
	// 	// runner.PrintResult(res.Items, os.Stdout)
	// case *oapi.HtcProjectsGetForbidden,
	// 	*oapi.HtcProjectsGetUnauthorized:
	// 	return nil, fmt.Errorf("forbidden: %s", res)
	// }

	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Submit(cmd *cobra.Command, args []string) error {
	projectId, err := cmd.Flags().GetString("project-id")
	if err != nil {
		return fmt.Errorf("Error setting project id: %w", err)
	}

	taskId, err := cmd.Flags().GetString("task-id")
	if err != nil {
		return fmt.Errorf("Error setting task id: %w", err)
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
		r, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("Error opening %s: %w", err)
		}
		defer r.Close()
	}

	req := submitRequest{
		params: oapi.HtcProjectsProjectIdTasksTaskIdJobsBatchPostParams{
			ProjectId: projectId,
			TaskId:    taskId,
			Group:     oapi.NewOptString(group),
		},
	}
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&req.batch); err != nil {
		return fmt.Errorf("Error parsing %s: %w", args[0], err)
	}

	runner, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}
	if err := runner.UpdateToken(time.Now()); err != nil {
		return err
	}

	// log.Printf("submit: projectId=%s taskId=%s", projectId, taskId)

	return nil
	ctx := context.Background()
	res, err := submit(ctx, runner.Client, &req)
	if err != nil {
		return fmt.Errorf("Error on job submission: %w", err)
	}
	log.Printf("res: %#v", res)
	return nil

	// var items []oapi.HTCProject
	// var pageIndex string
	// for {
	// 	res, err := getProjects(ctx, runner.Client, pageIndex)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	items = append(items, res.Items...)
	// 	if len(items) > limit {
	// 		items = items[:limit]
	// 		break
	// 	}

	// 	pageIndex = res.Next.Value.Query().Get("pageIndex")
	// 	if pageIndex == "" {
	// 		break
	// 	}
	// }
	// return runner.PrintResult(items, os.Stdout)
}

var SubmitCmd = &cobra.Command{
	Use:   "submit JOB_BATCH_YAML",
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

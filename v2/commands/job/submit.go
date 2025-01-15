package job

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
)

type submitRequest struct {
	batch  []oapi.HTCJobSubmitRequest
	params oapi.SubmitJobsParams
}

func submit(ctx context.Context, c oapi.JobInvoker, r *submitRequest) (*oapi.HTCJobSubmitRequests, error) {
	res, err := c.SubmitJobs(ctx, r.batch, r.params)
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HTCJobSubmitRequests:
		return res, nil
	case *oapi.HTCRequestError:
		return nil, fmt.Errorf("%s: %s", res.ErrorDescription.Value, res.Message.Value)
	}

	return nil, fmt.Errorf("Unknown response type: %T", res)
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

	if len(args) != 1 {
		return fmt.Errorf("Error: job yaml not provided")
	}
	f, err := common.OpenArg(args[0])
	if err != nil {
		return err
	}
	defer f.Close()

	req := submitRequest{
		params: oapi.SubmitJobsParams{
			ProjectId: p.ProjectId,
			TaskId:    p.TaskId,
			Group:     oapi.NewOptString(group),
		},
	}

	dec := json.NewDecoder(f)
	if err := dec.Decode(&req.batch); err != nil {
		return fmt.Errorf("Error parsing %s: %v", args[0], err)
	}

	ctx := context.Background()
	res, err := submit(ctx, runner.Client, &req)
	if err != nil {
		return fmt.Errorf("Error on job submission: %v", err)
	}
	return runner.PrintResult(res, os.Stdout)
}

var SubmitCmd = &cobra.Command{
	Use:   "submit JSON_FILE",
	Args:  cobra.ExactArgs(1),
	Short: "Submits jobs for a given task and project",
	Run:   common.WrapRunE(Submit),
}

func init() {
	// Prepare a sample JSON payload for our example.
	job := oapi.HTCJobSubmitRequest{
		JobName: oapi.NewOptString("a-rescale-htc-job"),
		HtcJobDefinition: oapi.HTCJobDefinition{
			ImageName:  "rescale-rsj-load-test-image_alpine_x86:latest",
			MaxVCpus:   oapi.NewOptInt32(1),
			MaxMemory:  oapi.NewOptInt32(128),
			MaxDiskGiB: oapi.NewOptInt32(1),
			Commands:   []string{"/bin/sh", "-c", "sleep 5m; echo all done"},
			Envs: []oapi.EnvPair{
				oapi.EnvPair{"INPUT_BUCKET", "htc-rescale-bucket"},
			},
			ExecTimeoutSeconds: oapi.NewOptInt32(3600),
			Priority:           oapi.NewOptJobPriority(oapi.JobPriorityONDEMANDPRIORITY),
		},
		BatchSize: oapi.NewOptInt32(1),
		Regions: []oapi.RescaleRegion{
			oapi.RescaleRegionGCPEUWEST2,
		},
		RetryStrategy: oapi.NewOptHTCRetryStrategy(
			oapi.HTCRetryStrategy{MaxRetries: oapi.NewOptInt32(1)},
		),
	}

	b, err := json.MarshalIndent(&job, "", "  ")
	if err != nil {
		panic("Unable to serialize `job create` JSON example: " + err.Error())
	}

	SubmitCmd.Flags().String("project-id", "", "HTC project ID (required)")
	SubmitCmd.Flags().String("task-id", "", "HTC task ID (required)")
	SubmitCmd.Flags().String("group", "", "Group")

	SubmitCmd.Long = SubmitCmd.Short + `
JSON_FILE is a path to a JSON file or - for stdin.`
	SubmitCmd.Example = fmt.Sprintf(`
htc job submit - <<'EOF'
  %s
EOF`, string(b))
}

package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
	"github.com/rescale-labs/htc-cli/v2/config"
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

	userProvidedWorkingDirectory, err := cmd.Flags().GetString("working-dir")
	if err != nil {
		return fmt.Errorf("Error setting working-directory: %w", err)
	}

	envMap := make(map[string]string)
	env, err := cmd.Flags().GetString("env")
	if err != nil {
		return fmt.Errorf("Error setting env: %w", err)
	}
	if len(env) > 0 {
		for _, token := range strings.Split(env, ",") {
			key, value, ok := strings.Cut(token, "=")
			if !ok {
				return config.UsageErrorf("env option has invalid format")
			}
			envMap[key] = value
		}
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
			Group:     oapi.OptString{group, group != ""},
		},
	}

	dec := json.NewDecoder(f)
	if err := dec.Decode(&req.batch); err != nil {
		return fmt.Errorf("Error parsing %s: %v", args[0], err)
	}

	// Patch job environment with envMap
	if len(envMap) > 0 {
		for i := range req.batch {
			for k, v := range envMap {
				req.batch[i].HtcJobDefinition.Envs = append(req.batch[i].HtcJobDefinition.Envs,
					oapi.EnvPair{Name: k, Value: v})
			}
		}
	}

	// Only execute if user actually opted-in for a flag
	if userProvidedWorkingDirectory != "" {
		workingDir, err := getWorkingDir(userProvidedWorkingDirectory)
		if err != nil {
			return err
		}

		// Set the workingDirectory and CFS experimental
		for i := range req.batch {
			req.batch[i].Experimental = oapi.OptExperimentalFields{
				Value: oapi.ExperimentalFields{
					KubernetesSwap:   oapi.OptBool{Set: false},
					CloudFileSystems: oapi.OptBool{Value: true, Set: true},
				},
				Set: true,
			}
			req.batch[i].HtcJobDefinition.WorkingDir = oapi.OptString{Value: *workingDir, Set: true}
		}
	}

	ctx := context.Background()
	res, err := submit(ctx, runner.Client, &req)
	if err != nil {
		return fmt.Errorf("Error on job submission: %v", err)
	}
	return runner.PrintResult(res, os.Stdout)
}

func getWorkingDir(userPassedDir string) (*string, error) {
	if !path.IsAbs(userPassedDir) {
		return nil, errors.New("only absolute paths are allowed when using working directory flag")
	}
	_, err := os.Stat(userPassedDir)
	if err != nil {
		slog.Warn("Warning: passed directory is not present on the system! Job will still be submitted with provided", "path", userPassedDir)
	}
	returnDirectory := userPassedDir
	if returnDirectory == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("cannot get current working directory %v", err)
		}
		returnDirectory = cwd
	}
	return &returnDirectory, nil
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
	SubmitCmd.Flags().StringP("env", "e", "", "Set job environment variables using comma-delimited KEY=VALUE pairs")
	// We wanted to be able to pass in just the -W flag which would mean CWD.
	// There was a choice to create 2 separate flags (boolean:-W and string:--working-dir), but I did not want 2 flags
	SubmitCmd.Flags().StringP("working-dir", "W", "", "Set current working directory (pwd) for a job commands execution. If no input given assumes current directory from which command is executed. Experimental feature")
	// There are 2 caveats with NoOptDefVal:
	// 1. Without NoOptDefVal "flag needs an argument" gets raised, eg `htc job submit job_def.json -W` is now allowed even with default value. See: https://github.com/spf13/cobra/issues/1756
	// 2. NoOptDefVal forces "=" assignment. https://github.com/spf13/pflag/issues/321
	SubmitCmd.Flag("working-dir").NoOptDefVal = "."

	SubmitCmd.Long = SubmitCmd.Short + `
JSON_FILE is a path to a JSON file or - for stdin.`
	SubmitCmd.Example = fmt.Sprintf(`
htc job submit - <<'EOF'
  %s
EOF`, string(b))
}

package job

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
	"github.com/rescale-labs/htc-cli/v2/config"
)

func logs(ctx context.Context, c oapi.JobInvoker, projectId, taskId, jobId, pageIndex string) (*oapi.HTCJobLogs, error) {
	res, err := c.GetLogs(ctx, oapi.GetLogsParams{
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
	case *oapi.HTCJobLogs:
		return res, nil
	case *oapi.GetLogsUnauthorized,
		*oapi.GetLogsForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %T", res)
}

func Logs(cmd *cobra.Command, args []string) error {
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

	flags := cmd.Flags()
	limit, err := flags.GetInt("limit")
	if err != nil {
		return config.UsageErrorf("Error setting limit: %w", err)
	}

	ctx := context.Background()
	var pageIndex string
	total := 0

	if _, err := fmt.Fprintf(os.Stdout, "%-38s %19s\n", "Timestamp", "Message"); err != nil {
		return fmt.Errorf("Error: unable to write header")
	}

	for {
		// write rows as the batches come in up to the limit
		if limit > 0 && limit <= total {
			break
		}

		res, err := logs(ctx, runner.Client, p.ProjectId, p.TaskId, jobId, pageIndex)
		if err != nil {
			return err
		}

		currentPage := len(res.Items)
		if total+currentPage >= limit {
			currentPage = limit - total
		}

		// only write up to the current page limit
		if limit > 0 {
			err = writeRows(res.Items[:currentPage], os.Stdout)
		} else {
			err = writeRows(res.Items, os.Stdout)
		}
		if err != nil {
			return err
		}

		pageIndex = res.Next.Value.Query().Get("pageIndex")
		if pageIndex == "" {
			break
		}
		total += len(res.Items)
	}
	return nil
}

func writeRows(rows []oapi.HTCLogEvent, w io.Writer) error {
	for _, row := range rows {
		timestamp := time.Time(row.Timestamp.Value).Format(time.DateTime)
		if _, err := fmt.Fprintf(w, "%-38s %19s\n", timestamp, row.Message.Value); err != nil {
			return err
		}
	}
	return nil
}

var LogsCmd = &cobra.Command{
	Use:   "logs [JOB_UUID]",
	Short: "Returns latest HTC job logs given a job ID.",
	Run:   common.WrapRunE(Logs),
	Args:  cobra.ExactArgs(1),
}

func init() {
	flags := LogsCmd.Flags()

	flags.IntP("limit", "l", 0, "Limit response to N items")
	flags.String("project-id", "", "HTC project ID")
	flags.String("task-id", "", "HTC task ID")
}

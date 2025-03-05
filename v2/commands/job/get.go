package job

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
	"github.com/rescale-labs/htc-cli/v2/config"
	"github.com/rescale-labs/htc-cli/v2/tabler"
)

func getJobs(ctx context.Context, c oapi.JobInvoker, params *oapi.GetJobsParams) (*oapi.HTCJobs, error) {
	res, err := c.GetJobs(ctx, *params)
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *oapi.HTCJobs:
		return res, nil
	case *oapi.GetJobsUnauthorized, *oapi.GetJobsForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func getJob(ctx context.Context, c oapi.JobInvoker, projectId, taskId, jobId string) (*oapi.HTCJob, error) {
	params := oapi.GetJobParams{
		ProjectId: projectId,
		TaskId:    taskId,
		JobId:     jobId,
	}
	log.Printf("getJob: %s %s %s", projectId, taskId, jobId)
	res, err := c.GetJob(ctx, params)
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *oapi.HTCJob:
		return res, nil
	case *oapi.GetJobUnauthorized, *oapi.GetJobForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Get(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true, RequireTaskId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	if len(args) == 0 {
		flags := cmd.Flags()

		limit, err := flags.GetInt("limit")
		if err != nil {
			return config.UsageErrorf("Error setting limit: %w", err)
		}

		latest, err := flags.GetInt("latest")
		if err != nil {
			return config.UsageErrorf("Error setting latest: %w", err)
		}

		if latest > 0 && limit > 0 {
			return config.UsageErrorf("Cannot use limit and latest in same command")
		}

		sortOrReversedPassed := flags.Changed("sort") || flags.Changed("reverse")
		if latest > 0 && sortOrReversedPassed {
			return config.UsageErrorf("Latest cannot be used together with sort or reverse options.")
		}

		group, err := flags.GetString("group")
		if err != nil {
			return config.UsageErrorf("Error setting group: %w", err)
		}

		reverse, err := flags.GetBool("reverse")
		if err != nil {
			return config.UsageErrorf("Error setting reverse: %w", err)
		}

		// --latest N == --limit N --reverse --sort created
		if latest > 0 {
			limit = latest
			reverse = true
			sort = common.SortCreated
		}

		var items []oapi.HTCJob
		params := oapi.GetJobsParams{
			ProjectId: p.ProjectId,
			TaskId:    p.TaskId,
			Group:     oapi.OptString{group, group != ""},
			PageSize:  oapi.NewOptInt32(common.PageSize),
			ViewType:  oapi.NewOptViewType(oapi.ViewTypeFULL),
		}

		if filterStatus != "" {
			params.Status = oapi.NewOptRescaleJobStatus(filterStatus.ToRescaleStatus())
		}

		for {
			res, err := getJobs(ctx, runner.Client, &params)
			if err != nil {
				return err
			}
			items = append(items, res.Items...)
			if !(sortOrReversedPassed || latest > 0) && limit > 0 && len(items) >= limit {
				items = items[:limit]
				break
			}

			params.PageIndex = oapi.NewOptString(
				res.Next.Value.Query().Get("pageIndex"))
			if params.PageIndex.Value == "" {
				break
			}
		}

		var sortFunc func(a, b oapi.HTCJob) int
		switch common.SortOrder(sort) {
		case "", common.SortCreated:
			sortFunc = func(a, b oapi.HTCJob) int {
				return time.Time(a.CreatedAt.Value).Compare(
					time.Time(b.CreatedAt.Value))
			}

		case common.SortCompleted:
			sortFunc = func(a, b oapi.HTCJob) int {
				return time.Time(a.CompletedAt.Value).Compare(
					time.Time(b.CompletedAt.Value))
			}

		case common.SortStatus:
			sortFunc = func(a, b oapi.HTCJob) int {
				ret := cmp.Compare(a.Status.Value, b.Status.Value)
				if ret == 0 {
					return time.Time(a.CreatedAt.Value).Compare(
						time.Time(b.CreatedAt.Value))
				}
				return ret
			}

		default:
			panic("Unrecognized sort option")
		}

		if reverse {
			oldFunc := sortFunc
			sortFunc = func(a, b oapi.HTCJob) int {
				return -1 * oldFunc(a, b)
			}
		}

		slices.SortFunc(items, sortFunc)
		if limit > 0 && len(items) >= limit {
			items = items[:limit]
		}
		return runner.PrintResult(tabler.HTCJobs(items), os.Stdout)
	}

	job, err := getJob(ctx, runner.Client, p.ProjectId, p.TaskId, args[0])
	if err != nil {
		return err
	}
	return runner.PrintResult((*tabler.HTCJob)(job), os.Stdout)
}

var GetCmd = &cobra.Command{
	Use:   "get [JOB_UUID]",
	Short: "Returns HTC jobs in a given task.",
	Run:   common.WrapRunE(Get),
	Args:  cobra.RangeArgs(0, 1),
}

var sort common.SortOrder

const sortDefault = common.SortCompleted

var filterStatus common.StatusFilter

const filterDefault = ""

func init() {
	flags := GetCmd.Flags()

	flags.IntP("limit", "l", 0, "Limit printed response to N items")
	flags.IntP("latest", "L", 0, "Limit printed response to latest N items. (Equivalent to & overrides --limits N --reverse --sort created).")
	flags.String("project-id", "", "HTC project ID")
	flags.String("task-id", "", "HTC task ID")
	flags.String("group", "", "HTC job batch group")
	flags.Var(&sort, "sort", fmt.Sprintf(
		"Sort job output (%s, default %q)",
		strings.Join([]string{
			string(common.SortCompleted),
			string(common.SortCreated),
			string(common.SortStatus),
		}, "|"),
		sortDefault))
	flags.Var(&filterStatus, "filter", fmt.Sprintf(
		"Filter output by job status (%s, default %q)",
		strings.Join(func() []string {
			var strs []string
			var rescaleJobStatus oapi.RescaleJobStatus
			allStatuses := rescaleJobStatus.AllValues()
			for _, s := range allStatuses {
				strs = append(strs, string(s))
			}
			return strs
		}(), "|"),
		filterDefault))
	flags.BoolP("reverse", "r", false, "Reverse sort order")
}

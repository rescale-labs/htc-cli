package project

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"github.com/rescale-labs/htc-cli/v2/common"
)

func getRetentionPolicy(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.HtcProjectsProjectIdTaskRetentionPolicyGet(ctx,
		oapi.HtcProjectsProjectIdTaskRetentionPolicyGetParams{ProjectId: p.ProjectId})
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.TaskRetentionPolicy:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.HtcProjectsProjectIdTaskRetentionPolicyGetForbidden,
		*oapi.HtcProjectsProjectIdTaskRetentionPolicyGetUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}

	return fmt.Errorf("Unknown response type: %s", res)
}

var RetentionPolicyCmd = &cobra.Command{
	Use:   "retention-policy",
	Short: "Commands for project retention policies",
}

var RetentionPolicyGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns retention policy for an HTC project",
	Run:   common.WrapRunE(getRetentionPolicy),
	Args:  cobra.ExactArgs(0),
}

func init() {
	RetentionPolicyCmd.AddCommand(RetentionPolicyGetCmd)
}

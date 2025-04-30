package project

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
	res, err := runner.Client.GetProjectTaskRetentionPolicy(ctx,
		oapi.GetProjectTaskRetentionPolicyParams{ProjectId: p.ProjectId})
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.TaskRetentionPolicy:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.GetProjectTaskRetentionPolicyForbidden,
		*oapi.GetProjectTaskRetentionPolicyUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}

	return fmt.Errorf("Unknown response type: %s", res)
}

func putRetentionPolicy(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	var policy oapi.TaskRetentionPolicy
	if err := common.DecodeFile(&policy, args[0]); err != nil {
		return err
	}
	
	ctx := context.Background()
	res, err := runner.Client.PutProjectTaskRetentionPolicy(ctx,
		oapi.NewOptTaskRetentionPolicy(policy),
		oapi.PutProjectTaskRetentionPolicyParams{
			ProjectId: p.ProjectId,
		},
	)
	if err != nil {
		return err
	}
	
	switch res := res.(type) {
	case *oapi.TaskRetentionPolicy:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.PutProjectTaskRetentionPolicyUnauthorized, *oapi.PutProjectTaskRetentionPolicyForbidden:
		return fmt.Errorf("forbidden: %s", res)
	}
	return nil

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

var RetentionPolicyApplyCmd = &cobra.Command{
	Use:	"apply JSON_FILE",
	Short: 	"Apply task retention policy to an HTC project.",
	Run:	common.WrapRunE(putRetentionPolicy),
	Args:	cobra.ExactArgs(1),
}

func init() {
	RetentionPolicyCmd.AddCommand(RetentionPolicyGetCmd)


	// example retention policy JSON payload
	policy := oapi.TaskRetentionPolicy{
		ArchiveAfter: 24, // hours
		DeleteAfter: 168, // hours
	}
	b, err := json.MarshalIndent(&policy, "", "  ")
	if err != nil {
		panic("Unable to serialize `retention-policy apply` JSON example: " + err.Error())
	}
	RetentionPolicyApplyCmd.Long = RetentionPolicyApplyCmd.Short + `
JSON_FILE is a path to a JSON file or - for stdin.`
	RetentionPolicyApplyCmd.Example = fmt.Sprintf(`
htc project retention-policy apply - <<'EOF'
  %s
EOF`, string(b))

	RetentionPolicyCmd.AddCommand(RetentionPolicyApplyCmd)
}

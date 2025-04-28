package workspace

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

func TaskRetentionPolicyGet(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireWorkspaceId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.GetTaskRetentionPolicy(ctx, oapi.GetTaskRetentionPolicyParams{
		WorkspaceId: p.WorkspaceId,
	})
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.WorkspaceTaskRetentionPolicy:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.GetTaskRetentionPolicyForbidden,
		*oapi.GetTaskRetentionPolicyUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}
	return fmt.Errorf("Unknown response type: %s", res)
}

func TaskRetentionPolicyPut(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireWorkspaceId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	f, err := common.OpenArg(args[0])
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	var policy oapi.WorkspaceTaskRetentionPolicy
	if err := dec.Decode(&policy); err != nil {
		return err
	}


	ctx := context.Background()
	res, err := runner.Client.PutTaskRetentionPolicy(ctx,
		oapi.NewOptWorkspaceTaskRetentionPolicy(policy),
		oapi.PutTaskRetentionPolicyParams{
			WorkspaceId: p.WorkspaceId,
		},
	)
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.WorkspaceTaskRetentionPolicy:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.PutTaskRetentionPolicyMethodNotAllowed:
		return fmt.Errorf("not allowed: %s", res)
	case *oapi.PutTaskRetentionPolicyUnauthorized, *oapi.PutTaskRetentionPolicyForbidden:
		return fmt.Errorf("forbidden: %s", res)
	}
	return nil
}

var RetentionPolicyCmd = &cobra.Command{
	Use: 	"retention-policy",
	Short: 	"Commands for workspace-scoped task retention policy",
}

var RetentionPolicyGetCmd = &cobra.Command{
	Use: 	"get",
	Short: 	"Returns task retention policy to a workspace.",
	Run:	common.WrapRunE(TaskRetentionPolicyGet),
	Args: 	cobra.ExactArgs(0),
}

var RetentionPolicyApplyCmd = &cobra.Command{
	Use:	"apply JSON_FILE",
	Short: 	"Apply task retention policy to a workspace.",
	Run:	common.WrapRunE(TaskRetentionPolicyPut),
	Args:	cobra.ExactArgs(1),
}

func init() {
	// example retention policy JSON payload

	policy := oapi.WorkspaceTaskRetentionPolicy{
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
htc workspace retention-policy apply - <<'EOF'
  %s
EOF`, string(b))

	RetentionPolicyCmd.AddCommand(RetentionPolicyGetCmd)
	RetentionPolicyCmd.AddCommand(RetentionPolicyApplyCmd)
}

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

func getLimits(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.GetProjectLimits(ctx,
		oapi.GetProjectLimitsParams{ProjectId: p.ProjectId})
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.HTCProjectLimits:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.GetProjectLimitsForbidden, *oapi.GetProjectLimitsUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}

	return fmt.Errorf("Unknown response type: %s", res)
}

func postLimits(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	var limits oapi.HTCLimitCreate
	if err := common.DecodeFile(&limits, args[0]); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.CreateProjectLimit(ctx,
		oapi.NewOptHTCLimitCreate(limits),
		oapi.CreateProjectLimitParams{
			ProjectId: p.ProjectId,
		},
	)
	if err != nil {
		return err
	}
	switch res := res.(type) {
	case *oapi.HTCProjectLimit:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.CreateProjectLimitUnauthorized, *oapi.CreateProjectLimitForbidden:
		return fmt.Errorf("forbidden: %s", res)
	}
	return nil
}

var LimitsCmd = &cobra.Command{
	Use:   "limits",
	Short: "Commands for project limits",
}

var LimitsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns limits for an HTC project",
	Run:   common.WrapRunE(getLimits),
	Args:  cobra.ExactArgs(0),
}

var LimitsApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Sets limits for an HTC project",
	Run:   common.WrapRunE(postLimits),
	Args:  cobra.ExactArgs(1),
}
//
// var LimitsDeleteCmd = &cobra.Command{
// 	Use:   "delete",
// 	Short: "Deletes limits for an HTC project",
// 	Run:   common.WrapRunE(LimitsGet),
// 	Args:  cobra.ExactArgs(0),
// }

func init() {
	LimitsCmd.AddCommand(LimitsGetCmd)

	// example limit JSON payload
	limit := oapi.HTCLimitCreate{
		ModifierRole: "PROJECT_ADMIN",
		VCPUs: 32,
	}
	b, err := json.MarshalIndent(&limit, "", "  ")
	if err != nil {
		panic("Unable to serialize `limits apply` JSON example: " + err.Error())
	}
	LimitsApplyCmd.Long = LimitsApplyCmd.Short + `
JSON_FILE is a path to a JSON file or - for stdin.`
	LimitsApplyCmd.Example = fmt.Sprintf(`
htc project limits apply - <<'EOF'
  %s
EOF`, string(b))

	LimitsCmd.AddCommand(LimitsApplyCmd)
}

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
	"github.com/rescale-labs/htc-cli/v2/tabler"
)

func getDimensions(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.GetProjectDimensions(ctx,
		oapi.GetProjectDimensionsParams{ProjectId: p.ProjectId})
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.HTCProjectDimensions:
		return runner.PrintResult((*tabler.ComputeEnvs)(res), os.Stdout)
	case *oapi.GetProjectDimensionsForbidden,
		*oapi.GetProjectDimensionsUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}

	return fmt.Errorf("Unknown response type: %s", res)
}

func putDimensions(cmd *cobra.Command, args[] string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	var limits oapi.HTCProjectDimensions
	if err := common.DecodeFile(&limits, args[0]); err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.CreateProjectDimensions(ctx,
		limits,
		oapi.CreateProjectDimensionsParams{
			ProjectId: p.ProjectId,
		},
	)
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.HTCProjectDimensions:
		return runner.PrintResult((*tabler.ComputeEnvs)(res), os.Stdout)
	case *oapi.CreateProjectDimensionsUnauthorized, *oapi.CreateProjectDimensionsForbidden:
		return fmt.Errorf("forbidden: %s", res)
	}
	return nil
}

var DimensionsCmd = &cobra.Command{
	Use:   "dimensions",
	Short: "Commands for project dimensions",
}

var DimensionsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns dimensions for an HTC project",
	Run:   common.WrapRunE(getDimensions),
	Args:  cobra.ExactArgs(0),
}

var DimensionsApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Sets dimensions for an HTC project",
	Run:   common.WrapRunE(putDimensions),
	Args:  cobra.ExactArgs(1),
}

// var DimensionsDeleteCmd = &cobra.Command{
// 	Use:   "delete",
// 	Short: "Deletes dimensions for an HTC project",
// 	Run:   common.WrapRunE(DimensionsGet),
// 	Args:  cobra.ExactArgs(0),
// }

func init() {
	DimensionsCmd.AddCommand(DimensionsGetCmd)
	
	// example dimensions JSON payload
	derived := oapi.HTCComputeEnvironmentDerived{
		Architecture: oapi.NewOptHTCComputeEnvironmentDerivedArchitecture(oapi.HTCComputeEnvironmentDerivedArchitectureAARCH64),
		VCPUs: oapi.NewOptInt32(16),
		Memory: oapi.NewOptFloat64(64),
		Swap: oapi.NewOptBool(false),
	}
	dimensions := oapi.HTCProjectDimensions{
		oapi.HTCComputeEnvironment{
			MachineType: oapi.NewOptString("n2d-standard-16"),
			Priority: oapi.NewOptHTCComputeEnvironmentPriority(oapi.HTCComputeEnvironmentPriorityONDEMANDPRIORITY),
			Region: oapi.NewOptHTCComputeEnvironmentRegion(oapi.HTCComputeEnvironmentRegion(oapi.HTCRegionAdminSettingsRegionGCPEUWEST4)),
			ComputeScalingPolicy: oapi.NewOptHTCComputeEnvironmentComputeScalingPolicy(oapi.HTCComputeEnvironmentComputeScalingPolicyOPTIMIZEUTILIZATION),
			Hyperthreading: oapi.NewOptBool(false),
			Derived: oapi.OptHTCComputeEnvironmentDerived(oapi.NewOptHTCComputeEnvironmentDerived(derived)),
		},
	}
	b, err := json.MarshalIndent(&dimensions, "", "  ")
	if err != nil {
		panic("Unable to serialize `dimensions apply` JSON example: " + err.Error())
	}
	DimensionsApplyCmd.Long = DimensionsApplyCmd.Short +`
	JSON_FILE is a path to a JSON file or - for stdin.`
	DimensionsApplyCmd.Example = fmt.Sprintf(`
htc project limits apply - <<'EOF'
  %s
EOF`, string(b))
	DimensionsCmd.AddCommand(DimensionsApplyCmd)
}

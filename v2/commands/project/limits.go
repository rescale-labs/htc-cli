package project

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/common"
)

func LimitsGet(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	projectId := args[0]

	ctx := context.Background()
	res, err := runner.Client.HtcProjectsProjectIdLimitsGet(ctx,
		oapi.HtcProjectsProjectIdLimitsGetParams{ProjectId: projectId})
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.HTCProjectLimits:
		return runner.PrintResult(res, os.Stdout)
	case *oapi.HtcProjectsProjectIdLimitsGetForbidden,
		*oapi.HtcProjectsProjectIdLimitsGetUnauthorized:
		return fmt.Errorf("forbidden: %s", res)
	}

	return fmt.Errorf("Unknown response type: %s", res)
}

var LimitsCmd = &cobra.Command{
	Use:   "limits",
	Short: "Commands for project limits",
}

var LimitsGetCmd = &cobra.Command{
	Use:   "get PROJECT_ID",
	Short: "Returns limits for an HTC project",
	Run:   common.WrapRunE(LimitsGet),
	Args:  cobra.ExactArgs(1),
}

var LimitsApplyCmd = &cobra.Command{
	Use:   "apply PROJECT_ID",
	Short: "Sets limits for an HTC project",
	Run:   common.WrapRunE(LimitsGet),
	Args:  cobra.ExactArgs(1),
}

var LimitsDeleteCmd = &cobra.Command{
	Use:   "delete PROJECT_ID",
	Short: "Deletes limits for an HTC project",
	Run:   common.WrapRunE(LimitsGet),
	Args:  cobra.ExactArgs(1),
}

func init() {
	LimitsCmd.AddCommand(LimitsGetCmd)
}

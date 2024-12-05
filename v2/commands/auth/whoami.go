package auth

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/common"
)

func WhoAmI(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}

	ctx := context.Background()
	res, err := runner.Client.WhoAmI(ctx)
	if err != nil {
		return err
	}

	switch res := res.(type) {
	case *oapi.WhoAmI:
		return runner.PrintResult(res, os.Stdout)
		// return nil
	case *oapi.OAuth2ErrorResponse:
		return fmt.Errorf("auth error: %s", res.GetError().Value)
	}

	return fmt.Errorf("Unknown response type: %s", res)
}

var WhoAmICmd = &cobra.Command{
	Use:   "whoami",
	Short: "Returns status of current user from HTC API.",
	Long:  "Returns status of current user from HTC API.\n\nUses Rescale API key and does not fetch or save a bearer token.",
	Run:   common.WrapRunE(WhoAmI),
}

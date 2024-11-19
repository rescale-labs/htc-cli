package metrics

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/common"
)

func getMetrics(ctx context.Context, c *oapi.Client) (io.Reader, error) {
	res, err := c.HtcMetricsGet(ctx, oapi.HtcMetricsGetParams{
		AcceptEncoding: []string{"text/plain"},
	})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HtcMetricsGetOK:
		return res.Data, nil
	case *oapi.HtcMetricsGetForbidden,
		*oapi.HtcMetricsGetUnauthorized:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Get(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	ctx := context.Background()
	r, err := getMetrics(ctx, runner.Client)
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, r)
	return err
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns HTC metrics in OpenMetrics / prometheus format",
	Run:   common.WrapRunE(Get),
}

// func init() {
// 	GetCmd.Flags().IntP("limit", "l", 0, "Limit response to N items")
// }

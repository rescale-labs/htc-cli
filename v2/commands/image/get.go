package image

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/common"
)

func getImages(ctx context.Context, c *oapi.Client, projectId string) (*oapi.HTCImages, error) {
	res, err := c.HtcProjectsProjectIdContainerRegistryImagesGet(ctx,
		oapi.HtcProjectsProjectIdContainerRegistryImagesGetParams{
			ProjectId: projectId,
		})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HTCImages:
		return res, nil
	case *oapi.HtcProjectsProjectIdContainerRegistryImagesGetUnauthorized:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func Get(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	images, err := getImages(ctx, runner.Client, p.ProjectId)
	if err != nil {
		return err
	}
	return runner.PrintResult(images, os.Stdout)
}

var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Returns HTC projects in a workspace",
	// Long:
	Run: common.WrapRunE(Get),
}

func init() {
	GetCmd.Flags().String("project-id", "", "HTC Project ID")
}

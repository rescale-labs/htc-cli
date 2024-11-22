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

func getImage(ctx context.Context, c *oapi.Client, projectId, imageName string) (*oapi.HTCImageStatus, error) {
	res, err := c.HtcProjectsProjectIdContainerRegistryImagesImageNameGet(ctx,
		oapi.HtcProjectsProjectIdContainerRegistryImagesImageNameGetParams{
			ProjectId: projectId,
			ImageName: imageName,
		})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HTCImageStatus:
		return res, nil
	case *oapi.HtcProjectsProjectIdContainerRegistryImagesImageNameGetUnauthorized,
		*oapi.HtcProjectsProjectIdContainerRegistryImagesImageNameGetForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	case *oapi.HtcProjectsProjectIdContainerRegistryImagesImageNameGetNotFound:
		return nil, fmt.Errorf("image not found")
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

	if len(args) == 0 {
		images, err := getImages(ctx, runner.Client, p.ProjectId)
		if err != nil {
			return err
		}
		return runner.PrintResult(images, os.Stdout)
	}

	image, err := getImage(ctx, runner.Client, p.ProjectId, args[0])
	if err != nil {
		return err
	}
	return runner.PrintResult(image, os.Stdout)
}

var GetCmd = &cobra.Command{
	Use:   "get [IMAGE_NAME:TAG]",
	Short: "Returns all container images in a given project or details on a single, tagged image.",
	Args:  cobra.RangeArgs(0, 1),
	Run:   common.WrapRunE(Get),
}

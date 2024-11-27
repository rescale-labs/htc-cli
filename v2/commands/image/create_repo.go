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

func createRepo(ctx context.Context, c oapi.ImageInvoker, projectId, repoName string) (*oapi.HTCRepository, error) {
	res, err := c.CreateRepo(ctx,
		oapi.CreateRepoParams{
			ProjectId: projectId,
			RepoName:  repoName,
		})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HTCRepository:
		return res, nil
	case *oapi.CreateRepoUnauthorized,
		*oapi.CreateRepoForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func CreateRepo(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	repo, err := createRepo(ctx, runner.Client, p.ProjectId, args[0])
	if err != nil {
		return err
	}
	return runner.PrintResult(repo, os.Stdout)
}

var CreateRepoCmd = &cobra.Command{
	Use:   "create-repo IMAGE_NAME",
	Short: "Creates a private container repository of a given name, within a given project.",
	// Long:
	Run:  common.WrapRunE(CreateRepo),
	Args: cobra.ExactArgs(1),
}

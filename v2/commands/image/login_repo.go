package image

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/rescale/htc-storage-cli/v2/common"
)

func LoginRepo(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	ctx := context.Background()
	_, err = login(ctx, runner.Client, p.ProjectId)
	return err
}

var LoginRepoCmd = &cobra.Command{
	Use:   "login-repo",
	Short: "Logs docker/podman into the given project's container registry",
	Run:   common.WrapRunE(LoginRepo),
}

package image

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
	"github.com/rescale/htc-storage-cli/v2/common"
)

func getToken(ctx context.Context, c *oapi.Client, projectId string) ([]byte, error) {
	res, err := c.HtcProjectsProjectIdContainerRegistryTokenGet(ctx,
		oapi.HtcProjectsProjectIdContainerRegistryTokenGetParams{
			ProjectId: projectId,
		})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HtcProjectsProjectIdContainerRegistryTokenGetOKHeaders:
		data, err := ioutil.ReadAll(res.Response.Data)
		if err != nil {
			return nil, err
		}
		return data, nil
	case *oapi.HtcProjectsProjectIdContainerRegistryTokenGetUnauthorized,
		*oapi.HtcProjectsProjectIdContainerRegistryTokenGetForbidden:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

func getLoginArgs(token []byte, registry string) ([]string, error) {
	// Default to "docker" but also try "podman"
	docker, err := exec.LookPath("docker")
	if err != nil {
		if docker, err = exec.LookPath("podman"); err != nil {
			return nil, err
		}
	}

	return []string{
		docker, "login",
		"--username", "AWS",
		"--password", string(token),
		registry,
	}, nil
}

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
	token, err := getToken(ctx, runner.Client, p.ProjectId)
	if err != nil {
		return err
	}

	images, err := getImages(ctx, runner.Client, p.ProjectId)
	if err != nil {
		return err
	}

	// Removing trailing / from registry; docker/podman won't accept it.
	registry := strings.TrimRight(images.ContainerRegistry.Value, "/")
	loginArgs, err := getLoginArgs(token, registry)
	if err != nil {
		return err
	}
	// log.Printf("login command: %s", strings.Join(loginArgs, " "))

	loginCmd := exec.CommandContext(ctx, loginArgs[0], loginArgs[1:]...)
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr
	return loginCmd.Run()
}

var LoginRepoCmd = &cobra.Command{
	Use:   "login-repo",
	Short: "Logs docker/podman into the given project's container registry",
	Run:   common.WrapRunE(LoginRepo),
}

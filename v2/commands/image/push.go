package image

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/rescale/htc-storage-cli/v2/common"
)

func tagImage(ctx context.Context, image, targetName string) error {
	docker, err := getDocker()
	if err != nil {
		return err
	}
	tagCmd := exec.CommandContext(
		ctx, docker, "tag", image, targetName)
	tagCmd.Stdout = os.Stdout
	tagCmd.Stderr = os.Stderr
	log.Printf(strings.Join(tagCmd.Args, " "))
	return tagCmd.Run()
}

func getPushArgs(image string) ([]string, error) {
	docker, err := getDocker()
	if err != nil {
		return nil, err
	}
	return []string{docker, "push", image}, nil
}

func push(cmd *cobra.Command, args []string) error {
	runner, err := common.NewRunnerWithToken(cmd, time.Now())
	if err != nil {
		return err
	}

	// Fetch repo name
	ctx := context.Background()
	p := common.IDParams{RequireProjectId: true}
	if err := runner.GetIds(&p); err != nil {
		return err
	}

	// Log in (even if we're already logged in, it's fast, and
	// there's no easy way to check our docker/podman login
	// programmatically.
	registry, err := login(ctx, runner.Client, p.ProjectId)
	if err != nil {
		return err
	}

	// Tag remote image from local image
	image := path.Join(registry, args[0])
	if err := tagImage(ctx, args[0], image); err != nil {
		return err
	}

	// Push!
	pushArgs, err := getPushArgs(image)
	if err != nil {
		return err
	}
	pushCmd := exec.CommandContext(ctx, pushArgs[0], pushArgs[1:]...)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	log.Printf(strings.Join(pushCmd.Args, " "))
	return pushCmd.Run()
}

var PushCmd = &cobra.Command{
	Use:   "push IMAGE_NAME:VERSION",
	Short: "Pushes container image to an HTC project's container registry",
	Run:   common.WrapRunE(push),
}

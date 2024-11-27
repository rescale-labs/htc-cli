package image

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
)

//
// HTC API
//

// Returns repo name and all images for a project from HTC API
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

// Returns repository name from HTCImages
func getRegistry(ctx context.Context, c *oapi.Client, projectId string) (string, error) {
	res, err := c.HtcProjectsProjectIdGet(ctx,
		oapi.HtcProjectsProjectIdGetParams{
			ProjectId: projectId,
		})
	if err != nil {
		return "", err
	}

	switch res := res.(type) {
	case *oapi.HTCProject:
		// Removing trailing / from registry; docker/podman won't accept it.
		return strings.TrimRight(res.ContainerRegistry.Value, "/"), nil
	case *oapi.HtcProjectsProjectIdGetUnauthorized, *oapi.HtcProjectsProjectIdGetForbidden:
		return "", fmt.Errorf("forbidden: %s", res)
	}
	return "", fmt.Errorf("Unknown response type: %s", res)

}

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

//
// Docker
//

var _docker string
var dockerMutex sync.Mutex

// Returns path to docker or podman, else an error.
func getDocker() (string, error) {
	// Cache the value, but safely.
	dockerMutex.Lock()
	defer dockerMutex.Unlock()
	if _docker != "" {
		return _docker, nil
	}
	docker, err := exec.LookPath("docker")
	if err != nil {
		if docker, err = exec.LookPath("podman"); err != nil {
			return "", err
		}
	}
	_docker = docker
	return _docker, nil
}

func getLoginArgs(registry string) ([]string, error) {
	docker, err := getDocker()
	if err != nil {
		return nil, err
	}
	return []string{
		docker, "login",
		"--username", "AWS",
		"--password-stdin",
		registry,
	}, nil
}

// Logs docker/podman into ECR registry. Returns registry name.
func login(ctx context.Context, c *oapi.Client, projectId string) (string, error) {
	token, err := getToken(ctx, c, projectId)
	if err != nil {
		return "", err
	}

	registry, err := getRegistry(ctx, c, projectId)
	if err != nil {
		return "", err
	}
	loginArgs, err := getLoginArgs(registry)
	if err != nil {
		return "", err
	}

	loginCmd := exec.CommandContext(ctx, loginArgs[0], loginArgs[1:]...)
	loginCmd.Stdin = bytes.NewReader(token)
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr
	if err := loginCmd.Run(); err != nil {
		return "", err
	}
	return registry, nil
}

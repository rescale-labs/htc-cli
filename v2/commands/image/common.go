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
func getImages(ctx context.Context, c oapi.ImageInvoker, projectId string) (*oapi.HTCImages, error) {
	res, err := c.GetImages(ctx,
		oapi.GetImagesParams{
			ProjectId: projectId,
		})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.HTCImages:
		return res, nil
	case *oapi.GetImagesUnauthorized:
		return nil, fmt.Errorf("forbidden: %s", res)
	}
	return nil, fmt.Errorf("Unknown response type: %s", res)
}

type imageProjectInvoker interface {
	oapi.ImageInvoker
	oapi.ProjectInvoker
}

// Returns repository name from HTCImages
func getRegistry(ctx context.Context, c imageProjectInvoker, projectId string) (string, error) {
	res, err := c.GetProject(ctx,
		oapi.GetProjectParams{
			ProjectId: projectId,
		})
	if err != nil {
		return "", err
	}

	switch res := res.(type) {
	case *oapi.HTCProject:
		// Removing trailing / from registry; docker/podman won't accept it.
		return strings.TrimRight(res.ContainerRegistry.Value, "/"), nil
	case *oapi.GetProjectUnauthorized, *oapi.GetProjectForbidden:
		return "", fmt.Errorf("forbidden: %s", res)
	}
	return "", fmt.Errorf("Unknown response type: %s", res)

}

func getToken(ctx context.Context, c oapi.ImageInvoker, projectId string) ([]byte, error) {
	res, err := c.GetRegistryToken(ctx,
		oapi.GetRegistryTokenParams{
			ProjectId: projectId,
		})
	if err != nil {
		return nil, err
	}

	switch res := res.(type) {
	case *oapi.GetRegistryTokenOKHeaders:
		data, err := ioutil.ReadAll(res.Response.Data)
		if err != nil {
			return nil, err
		}
		return data, nil
	case *oapi.GetRegistryTokenUnauthorized,
		*oapi.GetRegistryTokenForbidden:
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
func login(ctx context.Context, c imageProjectInvoker, projectId string) (string, error) {
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

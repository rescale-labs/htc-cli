package main

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"github.com/rescale/htc-storage-cli/cli"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"
)

func main() {
	ctx := context.Background()
	client, err := getGoogleClient(ctx)

	if err != nil {
		log.Fatalf("error creating client %s", err)
	}

	if len(os.Args) < 2 {
		err = errors.New("not enough arguments")
	}

	src, dest := cli.ParseArgs(os.Args[2:])

	if strings.HasPrefix(src, "gs://") {
		err = cli.Download(ctx, client, src, dest)
	} else {
		err = cli.Upload(ctx, client, src, dest)
	}

	if err != nil {
		log.Fatalf("error running command %s", err)
	}
}

func getGoogleClient(ctx context.Context) (*storage.Client, error) {
	credentialsString, err := getGoogleCredentials()
	if err != nil {
		return nil, err
	}

	credentials, err := google.CredentialsFromJSON(ctx, []byte(credentialsString), "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return nil, err
	}

	return storage.NewClient(ctx, option.WithCredentials(credentials))
}

func getGoogleCredentials() (string, error) {
	credentials, ok := os.LookupEnv("GCP_APPLICATION_CREDENTIALS")
	if !ok {
		return "", errors.New("error GCP_APPLICATION_CREDENTIALS not set")
	}

	return credentials, nil
}

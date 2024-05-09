package cli

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"os"
	"regexp"
	"strings"
)

func ParseBucket(bucketPath string) (string, string, error) {
	if !strings.HasPrefix(bucketPath, "gs://") {
		return "", "", errors.New("invalid bucket. Bucket must start with gs://")
	}

	re := regexp.MustCompile("^gs://(.+?)/(.+)")
	match := re.FindStringSubmatch(bucketPath)

	return match[1], strings.TrimRight(match[2], "/"), nil
}

func ParseArgs(args []string) TransferOptions {
	help := flag.Bool("h", false, "help message")
	cmd := flag.NewFlagSet("cp", flag.ContinueOnError)
	parallel := cmd.Int("p", 10, "Number of parallel transfers")

	err := cmd.Parse(args)

	if err != nil {
		return TransferOptions{"", "", 0, errors.New("error parsing args")}
	}

	if *help || len(cmd.Args()) != 2 {
		Usage()
	}

	src := cmd.Arg(0)
	dest := cmd.Arg(1)
	return TransferOptions{src, dest, *parallel, nil}
}

func Usage() {
	const usage = `Usage: htccli cp <src> <dst>


Available commands:
    cp      uploads or downloads one or more files to a destination URL

cp arguments:
    src		  path or cloud storage URI
    dest		path or cloud storage URI

Environment variables:
    GCP_APPLICATION_CREDENTIALS		JSON string containing GCP application credentials`
	fmt.Fprintf(os.Stderr, usage)

	os.Exit(1)
}

func GetGoogleClient(ctx context.Context) (*storage.Client, error) {
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
	credentials := os.Getenv("GCP_APPLICATION_CREDENTIALS")
	if credentials == "" {
		return "", errors.New("error GCP_APPLICATION_CREDENTIALS not set")
	}

	return credentials, nil
}

type TransferOptions struct {
	sourcePath      string
	destinationPath string
	parallelization int
	err             error
}

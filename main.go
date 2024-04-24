package main

import (
	"cloud.google.com/go/storage"
	"context"
	"flag"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	parseArgs()
}

func parseArgs() {
	validateArgs()

	ctx := context.Background()
	credentialsString := getGoogleCredentials()
	credentials, err := google.CredentialsFromJSON(ctx, []byte(credentialsString), "https://www.googleapis.com/auth/cloud-platform")

	if err != nil {
		log.Fatal("Error GCP_APPLICATION_CREDENTIALS not set")
		os.Exit(1)
	}

	client, err := storage.NewClient(ctx, option.WithCredentials(credentials))

	if err != nil {
		log.Fatalf("Error creating client %s", err)
		os.Exit(1)
	}

	if os.Args[1] == "upload" {
		upload(client, ctx)
	} else {
		download(client)
	}
}

func uploadFile(client *storage.Client, ctx context.Context, bucket string, object string, localFile string) {
	f, err := os.Open(localFile)
	if err != nil {
		log.Fatalf("os.Open: %w", err)
		os.Exit(1)
	}
	defer f.Close()

	o := client.Bucket(bucket).Object(object)
	o = o.If(storage.Conditions{DoesNotExist: true})

	writer := o.NewWriter(ctx)

	if _, err = io.Copy(writer, f); err != nil {
		log.Fatalf("io.Copy error: %w", err)
		return
	}

	if err := writer.Close(); err != nil {
		log.Fatalf("Writer.Close: %w", err)
		return
	}

	log.Printf("Blob %s uploaded.\n", object)
}

func getGoogleCredentials() string {
	credentials, ok := os.LookupEnv("GCP_APPLICATION_CREDENTIALS")
	if !ok {
		log.Fatal("Error GCP_APPLICATION_CREDENTIALS not set")
		os.Exit(1)
	}

	return credentials
}

func upload(client *storage.Client, ctx context.Context) {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	src := uploadCmd.String("src", "", "the source path to upload from")
	dest := uploadCmd.String("dest", "", "the destination bucket to upload to")
	region := uploadCmd.String("region", "", "the region to upload to")
	uploadCmd.Parse(os.Args[2:])

	bucket, path := parseBucket(*dest)

	object := fmt.Sprintf("%s/%s", path, "env")

	uploadFile(client, ctx, bucket, object, "env")

	//it := client.Bucket(bucket).Objects(ctx, &storage.Query{Prefix: path})
	//
	//for {
	//	attrs, err := it.Next()
	//	if err == iterator.Done {
	//		break
	//	}
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Printf("Name is %s", attrs.Name)
	//}

	log.Printf("Upload command: %s\n", uploadCmd.Args())
	log.Printf("src: %s\n", *src)
	log.Printf("dest: %s\n", *dest)
	log.Printf("region: %s\n", *region)
}

func parseBucket(bucketPath string) (string, string) {
	if !strings.HasPrefix(bucketPath, "gs://") {
		log.Fatal("Invalid bucket. Bucket must start with gs://")
		os.Exit(1)
	}

	re := regexp.MustCompile("gs://(.+?)/(.+)")
	match := re.FindStringSubmatch(bucketPath)

	return match[1], match[2]
}

func download(client *storage.Client) {
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	src := downloadCmd.String("src", "", "the source bucket to download from")
	dest := downloadCmd.String("dest", "", "the destination path to download to")
	region := downloadCmd.String("region", "", "the region to download from")
	downloadCmd.Parse(os.Args[2:])

	log.Printf("Download command: %s\n", downloadCmd.Args())
	log.Printf("src: %s\n", *src)
	log.Printf("dest: %s\n", *dest)
	log.Printf("region: %s\n", *region)
}

func validateArgs() {
	if len(os.Args) < 2 {
		log.Fatal("expected 'upload' or 'download' subcommands")
		os.Exit(1)
	}

	if os.Args[1] != "upload" && os.Args[1] != "download" {
		log.Fatal("expected one of the following sub commands 'upload' or 'download'")
		os.Exit(1)
	}
}

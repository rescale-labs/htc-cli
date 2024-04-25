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

type FileUpload struct {
	localPath  string
	remotePath string
}

type UploadResponse struct {
	remoteObject string
	success      bool
	error        error
}

func main() {
	validateArgs()

	ctx := context.Background()
	credentialsString := getGoogleCredentials()
	credentials, err := google.CredentialsFromJSON(ctx, []byte(credentialsString), "https://www.googleapis.com/auth/cloud-platform")

	if err != nil {
		log.Fatal("Error loading credentials")
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

func getGoogleCredentials() string {
	credentials, ok := os.LookupEnv("GCP_APPLICATION_CREDENTIALS")
	if !ok {
		log.Fatal("Error GCP_APPLICATION_CREDENTIALS not set")
		os.Exit(1)
	}

	return credentials
}

func uploadDirectory(client *storage.Client, ctx context.Context, bucket string, remotePath string, localPath string) {

	entries, err := os.ReadDir(localPath)
	if err != nil {
		log.Fatalf("Error opening path %s", localPath)
		os.Exit(1)
	}

	numJobs := 5
	filesToProcess := make(chan FileUpload, numJobs)
	uploadedResponses := make(chan UploadResponse, len(entries))
	responses := make([]UploadResponse, len(entries))

	// adding 5 workers for processing the queue
	for worker := 0; worker < 5; worker++ {
		go uploadWorker(client, ctx, bucket, filesToProcess, uploadedResponses)
	}

	// puts the files into the queue to process
	for _, e := range entries {
		remoteObject := fmt.Sprintf("%s/%s", remotePath, e.Name())
		localFile := fmt.Sprintf("%s/%s", localPath, e.Name())

		filesToProcess <- FileUpload{localFile, remoteObject}
	}

	// no more work to add to the queue
	close(filesToProcess)

	// reads from the results if file was uploaded properly
	for w := 0; w < len(entries); w++ {
		responses = append(responses, <-uploadedResponses)
	}

	failedUploads := strings.Builder{}
	failedUploads.WriteString("Failed to upload files [")
	for _, response := range responses {
		if response.success == false {
			failedUploads.WriteString(fmt.Sprintf("%s ", response.remoteObject))
		}
	}
	failedUploads.WriteString("]")

	log.Print(failedUploads.String())
}

func uploadWorker(client *storage.Client, ctx context.Context, bucket string, filesToProcess <-chan FileUpload, uploaded chan<- UploadResponse) {
	// processing the elements in the queue
	for localFile := range filesToProcess {
		uploaded <- uploadFile(client, ctx, bucket, localFile.remotePath, localFile.localPath)
	}
}

func uploadFile(client *storage.Client, ctx context.Context, bucket string, object string, localFile string) UploadResponse {
	f, err := os.Open(localFile)
	if err != nil {
		log.Fatal("Failed to open local file")
		return UploadResponse{object, false, err}
	}
	defer f.Close()

	o := client.Bucket(bucket).Object(object)

	writer := o.NewWriter(ctx)

	if _, err = io.Copy(writer, f); err != nil {
		log.Fatalf("io.Copy error: %s", err)
		return UploadResponse{object, false, err}
	}

	if err := writer.Close(); err != nil {
		log.Fatalf("Writer.Close: %s", err)
		return UploadResponse{object, false, err}
	}

	log.Printf("Blob %s uploaded.\n", object)
	return UploadResponse{o.ObjectName(), true, nil}
}

func upload(client *storage.Client, ctx context.Context) {
	uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
	src := uploadCmd.String("src", "", "the source path to upload from")
	dest := uploadCmd.String("dest", "", "the destination bucket to upload to")
	uploadCmd.Parse(os.Args[2:])

	bucket, path := parseBucket(*dest)

	filePtr, err := os.Stat(*src)
	if err != nil {
		log.Fatal("Unable to check if path is file or directory")
		os.Exit(1)
	}

	if filePtr.IsDir() {
		uploadDirectory(client, ctx, bucket, path, *src)
	}

	if filePtr.Mode().IsRegular() {
		uploadFile(client, ctx, bucket, fmt.Sprintf("%s/%s", path, filePtr.Name()), *src)
	}

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

	downloadCmd.Parse(os.Args[2:])

	log.Printf("Download command: %s\n", downloadCmd.Args())
	log.Printf("src: %s\n", *src)
	log.Printf("dest: %s\n", *dest)
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

package main

import (
	"cloud.google.com/go/storage"
	"context"
	"flag"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type FileUpload struct {
	localPath  string
	remotePath string
}

type FileDownload struct {
	destinationFile string
	remoteObject    string
}

type UploadResponse struct {
	remoteObject string
	success      bool
	error        error
}

type DownloadResponse struct {
	localFile string
	success   bool
	error     error
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
		download(client, ctx)
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
	failedUploads := strings.Builder{}
	failedUploads.WriteString("Failed to upload files [")

	filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			objectPath, _ := strings.CutPrefix(path, localPath)
			if strings.HasPrefix(objectPath, "/") {
				objectPath = strings.TrimLeft(objectPath, "/")
			}
			if strings.HasPrefix(remotePath, "/") {
				remotePath = strings.TrimLeft(remotePath, "/")
			}
			sourceFilePath := fmt.Sprintf("%s/%s", strings.TrimRight(localPath, "/"), objectPath)
			remoteFilePath := fmt.Sprintf("%s/%s", remotePath, objectPath)

			log.Printf("Uploading %s to %s", sourceFilePath, remoteFilePath)
			result := uploadFile(client, ctx, bucket, remoteFilePath, sourceFilePath)
			if !result.success {
				failedUploads.WriteString(fmt.Sprintf("%s ", result.remoteObject))
			}
		}
		return nil
	})
	failedUploads.WriteString("]")

	log.Print(failedUploads.String())
}

func uploadFile(client *storage.Client, ctx context.Context, bucket string, object string, localFile string) UploadResponse {
	f, err := os.Open(localFile)
	if err != nil {
		log.Fatal("Failed to open local file")
		return UploadResponse{object, false, err}
	}
	defer f.Close()

	o := client.Bucket(bucket).Object(object)

	workerCtx, cancel := context.WithTimeout(ctx, time.Hour*1)
	defer cancel()

	writer := o.NewWriter(workerCtx)

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

func downloadFile(client *storage.Client, ctx context.Context, bucket string, object string, localFile string) DownloadResponse {

	destinationDirectory := filepath.Dir(localFile)
	ensureDirectoryExists(destinationDirectory)

	workerCtx, cancel := context.WithTimeout(ctx, time.Hour*1)
	defer cancel()

	filePtr, err := os.Create(localFile)
	if err != nil {
		log.Fatal("Failed to open local file")
		return DownloadResponse{localFile, false, err}
	}

	rc, err := client.Bucket(bucket).Object(object).NewReader(workerCtx)
	if err != nil {
		log.Fatalf("Object(%q).NewReader: %w", object, err)
		return DownloadResponse{localFile, false, err}
	}
	defer rc.Close()

	if _, err := io.Copy(filePtr, rc); err != nil {
		log.Fatalf("io.Copy error: %s", err)
		return DownloadResponse{localFile, false, err}
	}

	if err = filePtr.Close(); err != nil {
		log.Fatalf("f.Close: %w", err)
		return DownloadResponse{localFile, false, err}
	}

	log.Printf("Blob %v downloaded to local file %v\n", object, localFile)
	return DownloadResponse{localFile, true, nil}
}

func download(client *storage.Client, ctx context.Context) {
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	src := downloadCmd.String("src", "", "the source bucket to download from")
	dest := downloadCmd.String("dest", "", "the destination path to download to")

	downloadCmd.Parse(os.Args[2:])

	bucket, path := parseBucket(*src)

	ensureDirectoryExists(*dest)
	listAndDownloadObjects(client, ctx, bucket, path, *dest)
}

func listAndDownloadObjects(client *storage.Client, ctx context.Context, bucket string, path string, destinationDir string) {

	listCtx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	failedDownloads := strings.Builder{}
	failedDownloads.WriteString("Failed to download files [")

	it := client.Bucket(bucket).Objects(listCtx, &storage.Query{Prefix: path})
	page := iterator.NewPager(it, 50, "")
	for {
		var remoteObjects []*storage.ObjectAttrs
		nextPageToken, err := page.NextPage(&remoteObjects)

		if err != nil {
			log.Fatalf("Error getting next page of objects %w", err)
		}

		for _, object := range remoteObjects {
			objectPath, _ := strings.CutPrefix(object.Name, path)
			if strings.HasPrefix(objectPath, "/") {
				objectPath = strings.TrimLeft(objectPath, "/")
			}
			destinationFilePath := fmt.Sprintf("%s/%s", strings.TrimRight(destinationDir, "/"), objectPath)
			log.Printf("Downloading %s to %s", object.Name, destinationFilePath)
			response := downloadFile(client, ctx, bucket, object.Name, destinationFilePath)
			if !response.success {
				failedDownloads.WriteString(fmt.Sprintf("%s ", response.localFile))
			}
		}

		if nextPageToken == "" {
			break
		}
	}
	failedDownloads.WriteString("]")

	log.Print(failedDownloads.String())
}

func ensureDirectoryExists(dirName string) {

	info, err := os.Stat(dirName)
	if err == nil && info.IsDir() {
		return
	}

	// owner has all permissions and the rest have read and execute permissions
	err = os.MkdirAll(dirName, 0755)
	if err != nil {
		log.Fatalf("Error Making directory %w", err)
	}
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

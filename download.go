package main

import (
	"cloud.google.com/go/storage"
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileDownload struct {
	destinationFile string
	remoteObject    string
}

type DownloadResponse struct {
	localFile string
	success   bool
	error     error
}

func Download(ctx context.Context, client *storage.Client) error {
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	src := downloadCmd.String("src", "", "the source bucket to download from")
	dest := downloadCmd.String("dest", "", "the destination path to download to")

	downloadCmd.Parse(os.Args[2:])

	bucket, path := ParseBucket(*src)

	err := ensureDirectoryExists(*dest)

	if err != nil {
		return err
	}

	return listAndDownloadObjects(client, ctx, bucket, path, *dest)
}

func listAndDownloadObjects(client *storage.Client, ctx context.Context, bucket string, path string, destinationDir string) error {

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
			return err
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
	return nil
}

func downloadFile(client *storage.Client, ctx context.Context, bucket string, object string, localFile string) DownloadResponse {

	destinationDirectory := filepath.Dir(localFile)
	err := ensureDirectoryExists(destinationDirectory)

	if err != nil {
		return DownloadResponse{localFile, false, err}
	}

	workerCtx, cancel := context.WithTimeout(ctx, time.Hour*1)
	defer cancel()

	filePtr, err := os.Create(localFile)
	if err != nil {
		log.Fatal("Failed to open local file")
		return DownloadResponse{localFile, false, err}
	}

	rc, err := client.Bucket(bucket).Object(object).NewReader(workerCtx)
	if err != nil {
		return DownloadResponse{localFile, false, err}
	}
	defer rc.Close()

	if _, err = io.Copy(filePtr, rc); err != nil {
		return DownloadResponse{localFile, false, err}
	}

	if err = filePtr.Close(); err != nil {
		return DownloadResponse{localFile, false, err}
	}

	log.Printf("Blob %v downloaded to local file %v\n", object, localFile)
	return DownloadResponse{localFile, true, nil}
}

func ensureDirectoryExists(dirName string) error {
	info, err := os.Stat(dirName)
	if err == nil && info.IsDir() {
		return nil
	}

	// owner has all permissions and the rest have read and execute permissions
	err = os.MkdirAll(dirName, 0755)
	if err != nil {
		return err
	}

	return nil
}

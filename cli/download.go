package cli

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Download(ctx context.Context, client *storage.Client, src string, dest string) error {

	bucket, path := ParseBucket(src)

	err := ensureDirectoryExists(dest)

	if err != nil {
		return err
	}

	return listAndDownloadObjects(ctx, client, bucket, path, dest)
}

func listAndDownloadObjects(ctx context.Context, client *storage.Client, bucket string, path string, destinationDir string) error {

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
			err = downloadFile(ctx, client, bucket, object.Name, destinationFilePath)
			if err != nil {
				failedDownloads.WriteString(fmt.Sprintf("%s ", destinationFilePath))
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

func downloadFile(ctx context.Context, client *storage.Client, bucket string, object string, localFile string) error {

	destinationDirectory := filepath.Dir(localFile)
	err := ensureDirectoryExists(destinationDirectory)

	if err != nil {
		return err
	}

	workerCtx, cancel := context.WithTimeout(ctx, time.Hour*1)
	defer cancel()

	filePtr, err := os.Create(localFile)
	if err != nil {
		return err
	}

	rc, err := client.Bucket(bucket).Object(object).NewReader(workerCtx)
	if err != nil {
		return err
	}
	defer rc.Close()

	if _, err = io.Copy(filePtr, rc); err != nil {
		return err
	}

	if err = filePtr.Close(); err != nil {
		return err
	}

	log.Printf("Blob %v downloaded to local file %v\n", object, localFile)
	return nil
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

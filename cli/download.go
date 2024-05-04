package cli

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Download iterates over remote objects from a csp to download into a destination directory
// An error is returned if there was a failure listing or downloading files
// The local destination path is created if it does not exist
func Download(ctx context.Context, client *storage.Client, src string, dest string) error {
	bucket, remotePath, err := ParseBucket(src)
	if err != nil {
		return err
	}

	destinationDirectory := filepath.Dir(dest)
	err = os.MkdirAll(destinationDirectory, 0755)
	if err != nil {
		return err
	}

	return downloadObjects(ctx, client, bucket, remotePath, dest)
}

func downloadObjects(ctx context.Context, client *storage.Client, bucket string, remotePath string, destinationDir string) error {

	listCtx, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()

	var failedDownloads []string

	it := client.Bucket(bucket).Objects(listCtx, &storage.Query{Prefix: remotePath})
	page := iterator.NewPager(it, 50, "")
	for {
		var remoteObjects []*storage.ObjectAttrs
		nextPageToken, err := page.NextPage(&remoteObjects)

		if err != nil {
			return err
		}

		for _, object := range remoteObjects {
			destinationFilePath := getLocalDestination(object.Name, remotePath, destinationDir)
			err = downloadFile(ctx, client, bucket, object.Name, destinationFilePath)
			if err != nil {
				failedDownloads = append(failedDownloads, object.Name)
			}
		}

		if nextPageToken == "" {
			break
		}
	}

	if len(failedDownloads) != 0 {
		pathNames := ""
		for _, sourcePath := range failedDownloads {
			pathNames += " " + sourcePath
		}
		return errors.New(fmt.Sprintf("The following files failed to download: %s", pathNames))
	}
	return nil
}

func downloadFile(ctx context.Context, client *storage.Client, bucket string, object string, localFile string) error {

	destinationDirectory := filepath.Dir(localFile)
	err := os.MkdirAll(destinationDirectory, 0755)
	if err != nil {
		return err
	}

	workerCtx, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()

	filePtr, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	rc, err := client.Bucket(bucket).Object(object).NewReader(workerCtx)
	if err != nil {
		return err
	}
	defer rc.Close()

	if _, err = io.Copy(filePtr, rc); err != nil {
		return err
	}

	log.Printf("Blob %v downloaded to local file %v\n", object, localFile)
	return nil
}

func getLocalDestination(objectName string, remotePath string, destinationDir string) string {
	objectPath := strings.TrimPrefix(objectName, remotePath)
	objectPath = strings.TrimPrefix(objectPath, "/")
	localDestination := strings.TrimPrefix(fmt.Sprintf("%s/%s", strings.TrimSuffix(destinationDir, "/"), objectPath), "/")
	if objectName == remotePath {
		_, file := filepath.Split(objectName)
		localDestination = path.Join(localDestination, file)
	}
	return localDestination
}

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
	"sync"
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

func downloadObjects(ctx context.Context, client *storage.Client, bucket, remotePath, destinationDir string) error {
	var failedDownloads []string

	jobs := make(chan TransferObject)
	results := make(chan TransferObject)
	wg := sync.WaitGroup{}

	const numWorkers = 10

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go downloadWorker(ctx, client, bucket, jobs, results, &wg)
	}

	go func() {
		listCtx, cancel := context.WithTimeout(ctx, time.Hour)
		defer cancel()
		it := client.Bucket(bucket).Objects(listCtx, &storage.Query{Prefix: remotePath})
		page := iterator.NewPager(it, 50, "")
		for {
			var remoteObjects []*storage.ObjectAttrs
			nextPageToken, err := page.NextPage(&remoteObjects)
			if err != nil {
				break
			}

			for _, object := range remoteObjects {
				destinationFilePath := getLocalDestination(object.Name, remotePath, destinationDir)
				jobs <- TransferObject{object.Name, destinationFilePath, nil}
			}

			if nextPageToken == "" {
				break
			}
		}
		// when there are no more files to download we close the jobs channel wait then close the results
		close(jobs)
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result.err != nil {
			failedDownloads = append(failedDownloads, result.source)
		}
	}

	if len(failedDownloads) != 0 {
		pathNames := ""
		for _, sourcePath := range failedDownloads {
			pathNames += "\n" + sourcePath
		}
		return errors.New(fmt.Sprintf("The following files failed to download: %s", pathNames))
	}
	return nil
}

func downloadFile(ctx context.Context, client *storage.Client, bucket string, object string, localFile string) TransferObject {
	result := TransferObject{localFile, object, nil}
	destinationDirectory := filepath.Dir(localFile)
	err := os.MkdirAll(destinationDirectory, 0755)
	if err != nil {
		result.err = err
		return result
	}

	workerCtx, cancel := context.WithTimeout(ctx, time.Hour)
	defer cancel()

	filePtr, err := os.Create(localFile)
	if err != nil {
		result.err = err
		return result
	}
	defer func() {
		err = filePtr.Close()
		if err != nil {
			result.err = err
		}
	}()

	rc, err := client.Bucket(bucket).Object(object).NewReader(workerCtx)
	if err != nil {
		result.err = err
		return result
	}
	defer func() {
		err = rc.Close()
		if err != nil {
			result.err = err
		}
	}()

	if _, err = io.Copy(filePtr, rc); err != nil {
		result.err = err
		return result
	}

	log.Printf("Blob %v downloaded to local file %v\n", object, localFile)
	return result
}

func getLocalDestination(objectName string, remotePath string, destinationDir string) string {
	objectPath := strings.TrimPrefix(objectName, remotePath)
	objectPath = strings.TrimPrefix(objectPath, "/")
	destinationDir = strings.TrimSuffix(destinationDir, "/")
	destinationPath := path.Join(destinationDir, objectPath)
	// this check is to support downloading a single remote object, not just a path
	if objectName == remotePath {
		_, file := filepath.Split(objectName)
		destinationPath = path.Join(destinationPath, file)
	}
	return destinationPath
}

func downloadWorker(ctx context.Context, client *storage.Client, bucket string, jobs <-chan TransferObject, results chan<- TransferObject, wg *sync.WaitGroup) {
	for job := range jobs {
		results <- downloadFile(ctx, client, bucket, job.source, job.destination)
	}
	wg.Done()
}

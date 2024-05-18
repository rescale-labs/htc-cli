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
func Download(ctx context.Context, client *storage.Client, transfer *Transfer) error {
	src := transfer.sources[0]
	dest := transfer.destination
	bucket, remotePath, err := ParseBucket(src)
	if err != nil {
		return err
	}

	destinationDirectory := filepath.Dir(dest)
	err = os.MkdirAll(destinationDirectory, 0755)
	if err != nil {
		return err
	}

	return downloadObjects(ctx, client, bucket, remotePath, dest, transfer)
}

func downloadObjects(ctx context.Context, client *storage.Client, bucket, remotePath, destinationDir string, transfer *Transfer) error {
	var failedDownloads []string

	jobs := make(chan TransferResult)
	results := make(chan TransferResult)
	pageError := make(chan error)
	wg := sync.WaitGroup{}

	numWorkers := transfer.parallelization

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
				pageError <- err
				break
			}

			for _, object := range remoteObjects {
				destinationFilePath := getLocalDestination(object.Name, remotePath, destinationDir)
				jobs <- TransferResult{object.Name, destinationFilePath, nil}
				// if we are only downloading a single file we can break out of this loop
				if object.Name == remotePath {
					break
				}
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
			log.Printf("Result = %v", result.err)
			failedDownloads = append(failedDownloads, result.source)
		}
	}

	close(pageError)
	if err := <-pageError; err != nil {
		return err
	}

	if len(failedDownloads) != 0 {
		pathNames := ""
		for _, sourcePath := range failedDownloads {
			pathNames += "\n" + sourcePath
		}
		return errors.New(fmt.Sprintf("The following files failed to download:\n%s", pathNames))
	}
	return nil
}

func downloadFile(ctx context.Context, client *storage.Client, bucket string, object string, localFile string) (result TransferResult) {
	result = TransferResult{localFile, object, nil}
	destinationDirectory := filepath.Dir(localFile)
	err := os.MkdirAll(destinationDirectory, 0755)
	if err != nil {
		result.err = err
		return result
	}

	// TODO: reduce default timeout and make it configurable
	workerCtx, cancel := context.WithTimeout(ctx, time.Minute*20)
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

	return result
}

func getLocalDestination(objectName string, remotePath string, destination string) string {
	// this is if we cut based on an object file
	objectPath := strings.TrimPrefix(objectName, remotePath[:strings.LastIndex(remotePath, "/")])
	// this is if we want to instead cut based on a directory
	if objectName != remotePath {
		objectPath = strings.TrimPrefix(objectName, remotePath)
	}
	objectPath = strings.TrimPrefix(objectPath, "/")
	destination = strings.TrimSuffix(destination, "/")
	destinationPath := path.Join(destination, objectPath)
	// this check is to support downloading a single remote object, not just a path
	if objectName == remotePath {
		destinationPath = destination
	}
	return destinationPath
}

func downloadWorker(ctx context.Context, client *storage.Client, bucket string, jobs <-chan TransferResult, results chan<- TransferResult, wg *sync.WaitGroup) {
	for job := range jobs {
		results <- downloadFile(ctx, client, bucket, job.source, job.destination)
	}
	wg.Done()
}

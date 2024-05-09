package cli

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Upload iterates over local paths to upload as a remote object
// An error is returned if there was a failure uploading
func Upload(ctx context.Context, client *storage.Client, transfer *Transfer) error {
	dest := transfer.destinationPath
	bucket, path, err := ParseBucket(dest)
	if err != nil {
		return err
	}
	return uploadFiles(ctx, client, bucket, path, transfer)
}

func uploadFiles(ctx context.Context, client *storage.Client, bucket, remotePath string, transfer *Transfer) error {
	var failedUploads []string

	jobs := make(chan TransferResult)
	results := make(chan TransferResult)
	walkError := make(chan error)
	wg := sync.WaitGroup{}

	numWorkers := transfer.parallelization

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go uploadWorker(ctx, client, bucket, jobs, results, &wg)
	}

	go func() {
		var err error = nil
		for _, source := range transfer.sourcePaths {
			if _, err := os.Stat(source); err != nil {
				continue
			}
			err = filepath.Walk(source, func(pathStr string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					// if we upload a file we just want to get the directory of that file
					dir := filepath.Dir(pathStr)
					// if we upload a directory
					stat, err := os.Stat(source)
					if stat.IsDir() {
						dir = source
					}
					if err != nil {
						return err
					}
					objectPath := strings.TrimPrefix(pathStr, dir)
					objectPath = strings.TrimPrefix(objectPath, "/")
					remotePath = strings.TrimPrefix(remotePath, "/")
					remoteFilePath := path.Join(remotePath, objectPath)

					upload := TransferResult{pathStr, remoteFilePath, nil}
					jobs <- upload
				}
				return err
			})
		}
		// when there are no more files to upload we close the jobs channel wait then close the results
		close(jobs)
		wg.Wait()
		close(results)
		walkError <- err
	}()

	for result := range results {
		if result.err != nil {
			failedUploads = append(failedUploads, result.source)
		}
	}

	close(walkError)
	if err := <-walkError; err != nil {
		return err
	}

	if len(failedUploads) != 0 {
		pathNames := ""
		for _, sourcePath := range failedUploads {
			pathNames += "\n" + sourcePath
		}
		return errors.New(fmt.Sprintf("The following files failed to upload: %s", pathNames))
	}
	return nil
}

func uploadFile(ctx context.Context, client *storage.Client, bucket string, object string, localFile string) (result TransferResult) {
	result = TransferResult{localFile, object, nil}
	f, err := os.Open(localFile)
	if err != nil {
		result.err = err
		return result
	}

	defer func() {
		err = f.Close()
		if err != nil {
			result.err = err
		}
	}()

	o := client.Bucket(bucket).Object(object)

	// TODO: reduce default timeout and make it configurable
	workerCtx, cancel := context.WithTimeout(ctx, time.Minute*20)
	defer cancel()

	writer := o.NewWriter(workerCtx)
	defer func() {
		err = writer.Close()
		if err != nil {
			result.err = err
		}
	}()

	if _, err = io.Copy(writer, f); err != nil {
		result.err = err
		return result
	}

	log.Printf("Blob %s uploaded.\n", object)
	return result
}

func uploadWorker(ctx context.Context, client *storage.Client, bucket string, jobs <-chan TransferResult, results chan<- TransferResult, wg *sync.WaitGroup) {
	for job := range jobs {
		results <- uploadFile(ctx, client, bucket, job.destination, job.source)
	}
	wg.Done()
}

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
	dest := transfer.destination
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
	wg := sync.WaitGroup{}

	numWorkers := transfer.parallelization

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go uploadWorker(ctx, client, bucket, jobs, results, &wg)
	}

	go func() {
		for _, source := range transfer.sources {
			// TODO: Improve handling of walk errors and stat errors
			if _, err := os.Stat(source); err != nil {
				continue
			}
			filepath.Walk(source, func(pathStr string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					objectPath := strings.TrimPrefix(pathStr, "/")
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
	}()

	for result := range results {
		if result.err != nil {
			failedUploads = append(failedUploads, result.source)
		}
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

	log.Printf("Blob %s uploaded to %s", localFile, object)
	return result
}

func uploadWorker(ctx context.Context, client *storage.Client, bucket string, jobs <-chan TransferResult, results chan<- TransferResult, wg *sync.WaitGroup) {
	for job := range jobs {
		results <- uploadFile(ctx, client, bucket, job.destination, job.source)
	}
	wg.Done()
}

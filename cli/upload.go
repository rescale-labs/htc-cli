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
func Upload(ctx context.Context, client *storage.Client, transferOptions *TransferOptions) error {
	src := transferOptions.sourcePath
	dest := transferOptions.destinationPath
	bucket, path, err := ParseBucket(dest)
	if err != nil {
		return err
	}

	stat, err := os.Stat(src)
	if err != nil {
		return errors.New("unable to check if path is file or directory")
	}

	if stat.IsDir() {
		return uploadDirectory(ctx, client, bucket, path, src, transferOptions)
	} else if stat.Mode().IsRegular() {
		upload := uploadFile(ctx, client, bucket, fmt.Sprintf("%s/%s", path, stat.Name()), src)
		return upload.err
	} else {
		return fmt.Errorf("path %s is not a directory or file", src)
	}
}

func uploadDirectory(ctx context.Context, client *storage.Client, bucket, remotePath, localPath string, options *TransferOptions) error {
	var failedUploads []string

	jobs := make(chan TransferResult)
	results := make(chan TransferResult)
	walkError := make(chan error)
	wg := sync.WaitGroup{}

	numWorkers := options.parallelization
	log.Printf("%d", numWorkers)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go uploadWorker(ctx, client, bucket, jobs, results, &wg)
	}

	go func() {
		err := filepath.Walk(localPath, func(pathStr string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				objectPath := strings.TrimPrefix(pathStr, localPath)
				objectPath = strings.TrimPrefix(objectPath, "/")
				remotePath = strings.TrimPrefix(remotePath, "/")
				sourceFilePath := path.Join(strings.TrimSuffix(localPath, "/"), objectPath)
				remoteFilePath := path.Join(remotePath, objectPath)

				upload := TransferResult{sourceFilePath, remoteFilePath, nil}
				jobs <- upload
			}
			return err
		})
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
	workerCtx, cancel := context.WithTimeout(ctx, time.Hour*1)
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

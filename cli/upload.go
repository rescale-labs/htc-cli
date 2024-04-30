package cli

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Upload(ctx context.Context, client *storage.Client, src string, dest string) error {
	bucket, path := ParseBucket(dest)

	filePtr, err := os.Stat(src)
	if err != nil {
		return errors.New("unable to check if path is file or directory")
	}

	if filePtr.IsDir() {
		return uploadDirectory(client, ctx, bucket, path, src)
	} else if filePtr.Mode().IsRegular() {
		err = uploadFile(client, ctx, bucket, fmt.Sprintf("%s/%s", path, filePtr.Name()), src)
		return err
	} else {
		return errors.New("file pointer is not a directory or file")
	}
}

func uploadDirectory(client *storage.Client, ctx context.Context, bucket string, remotePath string, localPath string) error {
	failedUploads := strings.Builder{}
	failedUploads.WriteString("Failed to upload files [")

	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
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
			err = uploadFile(client, ctx, bucket, remoteFilePath, sourceFilePath)
			if err != nil {
				failedUploads.WriteString(fmt.Sprintf("%s ", remoteFilePath))
			}
		}
		return nil
	})
	failedUploads.WriteString("]")

	log.Print(failedUploads.String())
	return err
}

func uploadFile(client *storage.Client, ctx context.Context, bucket string, object string, localFile string) error {
	f, err := os.Open(localFile)
	if err != nil {
		log.Fatal("Failed to open local file")
		return err
	}
	defer f.Close()

	o := client.Bucket(bucket).Object(object)

	workerCtx, cancel := context.WithTimeout(ctx, time.Hour*1)
	defer cancel()

	writer := o.NewWriter(workerCtx)

	if _, err = io.Copy(writer, f); err != nil {
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}

	log.Printf("Blob %s uploaded.\n", object)
	return nil
}

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
	"time"
)

func Upload(ctx context.Context, client *storage.Client, src string, dest string) error {
	bucket, path, err := ParseBucket(dest)
	if err != nil {
		return err
	}

	stat, err := os.Stat(src)
	if err != nil {
		return errors.New("unable to check if path is file or directory")
	}

	if stat.IsDir() {
		return uploadDirectory(ctx, client, bucket, path, src)
	} else if stat.Mode().IsRegular() {
		err = uploadFile(ctx, client, bucket, fmt.Sprintf("%s/%s", path, stat.Name()), src)
		return err
	} else {
		return fmt.Errorf("path %s is not a directory or file", src)
	}
}

func uploadDirectory(ctx context.Context, client *storage.Client, bucket string, remotePath string, localPath string) error {
	var failedUploads []string

	err := filepath.Walk(localPath, func(pathStr string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			objectPath := strings.TrimPrefix(pathStr, localPath)
			objectPath = strings.TrimPrefix(objectPath, "/")
			remotePath = strings.TrimPrefix(remotePath, "/")

			sourceFilePath := path.Join(strings.TrimSuffix(localPath, "/"), objectPath)
			remoteFilePath := path.Join(remotePath, objectPath)

			log.Printf("Uploading %s to %s", sourceFilePath, remoteFilePath)
			err = uploadFile(ctx, client, bucket, remoteFilePath, sourceFilePath)
			if err != nil {
				failedUploads = append(failedUploads, sourceFilePath)
			}
		}
		return nil
	})

	if len(failedUploads) != 0 {
		pathNames := ""
		for _, sourcePath := range failedUploads {
			pathNames += " " + sourcePath
		}
		return errors.New(fmt.Sprintf("The following files failed to upload: %s", pathNames))
	}
	return err
}

func uploadFile(ctx context.Context, client *storage.Client, bucket string, object string, localFile string) error {
	f, err := os.Open(localFile)
	if err != nil {
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

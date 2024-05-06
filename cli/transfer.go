package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type TransferPath struct {
	src  string
	dest string
}

func NewTransfer(src string, dest string) *TransferPath {
	return &TransferPath{src, dest}
}

type TransferFile interface {
	Transfer(ctx context.Context) error
}

func (transfer *TransferPath) Transfer(ctx context.Context) error {
	client, err := GetGoogleClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	if strings.HasPrefix(transfer.src, "gs://") {
		err = Download(ctx, client, transfer.src, transfer.dest)
	} else if strings.HasPrefix(transfer.dest, "gs://") {
		err = Upload(ctx, client, transfer.src, transfer.dest)
	} else {
		err = localTransfer(transfer.src, transfer.dest)
	}
	return err
}

func localTransfer(src string, dest string) error {
	var failedCopies []string

	stat, err := os.Stat(src)
	if stat.IsDir() {
		err = filepath.Walk(src, func(pathStr string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				srcPath, _ := strings.CutPrefix(pathStr, src)
				destPath := path.Join(dest, srcPath)
				err = copyFile(pathStr, destPath)
				if err != nil {
					failedCopies = append(failedCopies, srcPath)
				}
			}
			return nil
		})

		if len(failedCopies) != 0 {
			pathNames := ""
			for _, sourcePath := range failedCopies {
				pathNames += " " + sourcePath
			}
			return errors.New(fmt.Sprintf("The following files failed to copy: %s", pathNames))
		}
	} else {
		copyFile(src, dest)
	}
	return err
}

func copyFile(src, dest string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destinationDirectory := filepath.Dir(dest)
	err = os.MkdirAll(destinationDirectory, 0755)
	if err != nil {
		return err
	}

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	log.Printf("File %s copied to %s.\n", src, dest)
	return err
}

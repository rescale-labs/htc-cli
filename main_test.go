package main

import (
	"context"
	"github.com/rescale/htc-storage-cli/cli"
	"os"
	"testing"
)

func TestParseBucket(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedBucket string
		expectedPath   string
	}{
		"Test Valid Bucket": {
			input:          "gs://my-bucket/object/prefix",
			expectedBucket: "my-bucket",
			expectedPath:   "object/prefix",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualBucket, actualPath, err := cli.ParseBucket(test.input)

			if err != nil {
				t.Errorf("There was an error parsing the bucket")
			}

			if actualBucket != test.expectedBucket {
				t.Errorf("Actual bucket %s did not equal expected bucket %s", actualBucket, test.expectedBucket)
			}

			if actualPath != test.expectedPath {
				t.Errorf("Actual path %s did not equal expected path %s", actualPath, test.expectedPath)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tests := map[string]struct {
		input        []string
		expectedSrc  string
		expectedDest string
	}{
		"Test Valid Args": {
			input:        []string{"cp", "gs://my-bucket/object/prefix", "local/path"},
			expectedSrc:  "gs://my-bucket/object/prefix",
			expectedDest: "local/path",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			src, dest := cli.ParseArgs(test.input[1:])
			if src != test.expectedSrc {
				t.Errorf("Actual src %s did not equal expected src %s", src, test.expectedSrc)
			}
			if dest != test.expectedDest {
				t.Errorf("Actual dest %s did not equal expected dest %s", dest, test.expectedDest)
			}
		})
	}
}

func TestTransfer(t *testing.T) {
	// create a test file
	localSrc := "test.txt"
	touchFile(localSrc)
	localDest := "test/test.txt"

	tests := map[string]struct {
		src  string
		dest string
	}{
		"Test Local Transfer": {
			src:  localSrc,
			dest: localDest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			transfer := cli.NewTransfer(test.src, test.dest)
			err := transfer.Transfer(context.Background())
			if err != nil {
				t.Errorf("Error transfering file %s", err)
			}

			stat, err := os.Stat(test.dest)
			if err != nil {
				t.Errorf("Error opening stat %s", err)
			}
			if !stat.Mode().IsRegular() {
				t.Errorf("Error destination file does not exist")
			}
		})
	}
	cleanUpFiles(localSrc, localDest)
}

func cleanUpFiles(src string, dest string) {
	os.Remove(src)
	os.Remove(dest)
}

func touchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

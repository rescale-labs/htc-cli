package main

import (
	"context"
	"errors"
	"github.com/rescale/htc-storage-cli/cli"
	"os"
	"testing"
)

func TestParseBucket(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedBucket string
		expectedPath   string
		expectedErr    error
	}{
		"Test Valid Bucket": {
			input:          "gs://my-bucket/object/prefix",
			expectedBucket: "my-bucket",
			expectedPath:   "object/prefix",
			expectedErr:    nil,
		},
		"Test Invalid Bucket": {
			input:          "my-bucket/object/prefix",
			expectedBucket: "",
			expectedPath:   "",
			expectedErr:    errors.New("invalid bucket. Bucket must start with gs://"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualBucket, actualPath, err := cli.ParseBucket(test.input)
			if err != nil && err.Error() != test.expectedErr.Error() {
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
	localSrc := "src"
	localDest := "test/test.txt"
	src, err := createTempFile("", localSrc)
	if err != nil {
		t.Errorf("Unable to create source temp file")
	}
	defer os.Remove(src.Name())

	tests := map[string]struct {
		src  string
		dest string
	}{
		"Test Local Transfer": {
			src:  src.Name(),
			dest: localDest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			transfer := cli.NewTransfer(test.src, test.dest)
			err = transfer.Transfer(context.Background())
			if err != nil {
				t.Errorf("Error transfering file %s", err)
			}
			defer os.Remove(test.dest)

			stat, err := os.Stat(test.dest)
			if err != nil {
				t.Errorf("Error opening stat %s", err)
			}
			if !stat.Mode().IsRegular() {
				t.Errorf("Error destination file does not exist")
			}
		})
	}
}

func createTempFile(dir string, name string) (*os.File, error) {
	file, err := os.CreateTemp(dir, name)
	if err != nil {
		return file, err
	}
	return file, err
}

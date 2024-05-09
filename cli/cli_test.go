package cli

import (
	"context"
	"os"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := map[string]struct {
		input        []string
		expectedSrc  string
		expectedDest string
		expectedErr  error
	}{
		"Test Valid Args": {
			input:        []string{"cp", "gs://my-bucket/object/prefix", "local/path"},
			expectedSrc:  "gs://my-bucket/object/prefix",
			expectedDest: "local/path",
			expectedErr:  nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			transferOpts, err := ParseArgs(test.input[1:])
			if err != test.expectedErr {
				t.Errorf("Actual error did not equal expected error")
			}
			if transferOpts.sourcePaths[0] != test.expectedSrc {
				t.Errorf("Actual src %s did not equal expected src %s", transferOpts.sourcePaths[0], test.expectedSrc)
			}
			if transferOpts.destinationPath != test.expectedDest {
				t.Errorf("Actual dest %s did not equal expected dest %s", transferOpts.destinationPath, test.expectedDest)
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
			transfer := &Transfer{[]string{test.src}, test.dest, 10}
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

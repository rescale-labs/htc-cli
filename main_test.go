package main

import (
	"github.com/rescale/htc-storage-cli/cli"
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
			actualBucket, actualPath := cli.ParseBucket(test.input)

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

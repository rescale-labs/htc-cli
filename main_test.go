package main

import (
	"errors"
	"github.com/rescale/htc-storage-cli/cli"
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

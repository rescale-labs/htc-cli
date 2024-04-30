package main

import (
	"cli"
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

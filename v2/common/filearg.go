package common

import (
	"io"
	"os"
)

// Returns an io.ReadCloser given an arg from the command line
// that is either a path to a file, or '-' for stdin.
func OpenArg(path string) (io.ReadCloser, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

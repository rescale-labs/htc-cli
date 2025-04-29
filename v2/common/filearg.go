package common

import (
	"encoding/json"
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

// Decode a provided JSON file into the template type
func DecodeFile(target any, path string) error {
    f, err := OpenArg(path)
    if err != nil {
        return err
    }
    defer f.Close()

    return json.NewDecoder(f).Decode(target)
}

package config

import "fmt"

// Custom error type used to differentiate between runtime errors and
// errors in argument/parameters passed to a command.
//
// Currently here in config, because config is where many common
// parameters are validated.
type UsageError struct {
	error
}

func UsageErrorf(msg string, args ...interface{}) error {
	return &UsageError{fmt.Errorf(msg, args...)}
}

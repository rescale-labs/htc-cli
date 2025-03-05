package common

import (
	"fmt"
	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"strings"
)

// SortOrder represents a sorting option for a command-line flag.
// It implements the [pflag.Value] interface, allowing it to be used as a flag value.
// See: https://pkg.go.dev/github.com/spf13/pflag#Value
//
// Supported values:
//   - "completed": Sort by completion time.
//   - "created": Sort by creation time.
//   - "status": Sort by status.
//
// Example usage with pflag:
//
//	var order SortOrder
//	flag.Var(&order, "sort", "Sorting order (completed|created|status)")
type SortOrder string

const (
	SortCompleted SortOrder = "completed"
	SortCreated   SortOrder = "created"
	SortStatus    SortOrder = "status"
)

func (s *SortOrder) String() string {
	return string(*s)
}

func (s *SortOrder) Set(v string) error {
	switch SortOrder(v) {
	case SortCompleted, SortCreated, SortStatus:
		*s = SortOrder(v)
		return nil
	default:
		return fmt.Errorf("%q is not a valid sort option", v)
	}
}

func (s *SortOrder) Type() string {
	return "string"
}

// StatusFilter is a filter type that represents a Rescale job status.
// It implements the [pflag.Value] interface, allowing it to be used as a flag value.
// See: https://pkg.go.dev/github.com/spf13/pflag#Value
//
// Supported values (case-insensitive):
//   - "SUBMITTEDTORESCALE": Job submitted to Rescale.
//   - "SUBMITTEDTOPROVIDER": Job submitted to the provider.
//   - "RUNNABLE": Job is ready to run.
//   - "STARTING": Job is starting.
//   - "RUNNING": Job is currently running.
//   - "SUCCEEDED": Job has successfully completed.
//   - "FAILED": Job has failed.
//
// Example usage with pflag:
//
//	var filter StatusFilter
//	flag.Var(&filter, "status", "Filter by job status (e.g., SUBMITTEDTORESCALE, RUNNING, SUCCEEDED)")
type StatusFilter oapi.RescaleJobStatus

func (sf *StatusFilter) String() string {
	return string(*sf)
}

// Set assigns a valid job status value to the StatusFilter.
// The input is case-insensitive.
func (s *StatusFilter) Set(val string) error {
	val = strings.ToUpper(val) // Case-insensitive
	switch val {
	case string(oapi.RescaleJobStatusSUBMITTEDTORESCALE):
		*s = StatusFilter(oapi.RescaleJobStatusSUBMITTEDTORESCALE)
		return nil
	case string(oapi.RescaleJobStatusSUBMITTEDTOPROVIDER):
		*s = StatusFilter(oapi.RescaleJobStatusSUBMITTEDTOPROVIDER)
		return nil
	case string(oapi.RescaleJobStatusRUNNABLE):
		*s = StatusFilter(oapi.RescaleJobStatusRUNNABLE)
		return nil
	case string(oapi.RescaleJobStatusSTARTING):
		*s = StatusFilter(oapi.RescaleJobStatusSTARTING)
		return nil
	case string(oapi.RescaleJobStatusRUNNING):
		*s = StatusFilter(oapi.RescaleJobStatusRUNNING)
		return nil
	case string(oapi.RescaleJobStatusSUCCEEDED):
		*s = StatusFilter(oapi.RescaleJobStatusSUCCEEDED)
		return nil
	case string(oapi.RescaleJobStatusFAILED):
		*s = StatusFilter(oapi.RescaleJobStatusFAILED)
		return nil
	default:
		panic("Unknown Rescale job status for filtering!")
	}
}

func (sf *StatusFilter) Type() string {
	return "StatusFilter"
}

// ToRescaleStatus converts a StatusFilter to an oapi.RescaleJobStatus.
func (sf StatusFilter) ToRescaleStatus() oapi.RescaleJobStatus {
	return oapi.RescaleJobStatus(sf)
}

// ToStatusFilter converts an oapi.RescaleJobStatus to a StatusFilter.
func ToStatusFilter(rjs oapi.RescaleJobStatus) StatusFilter { return StatusFilter(rjs) }

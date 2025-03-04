package common

import (
	"fmt"
	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
	"strings"
)

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

type StatusFilter oapi.RescaleJobStatus

func (sf *StatusFilter) String() string {
	return string(*sf)
}

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

func (sf StatusFilter) ToRescaleStatus() oapi.RescaleJobStatus {
	return oapi.RescaleJobStatus(sf)
}

func ToStatusFilter(rjs oapi.RescaleJobStatus) StatusFilter {
	return StatusFilter(rjs)
}

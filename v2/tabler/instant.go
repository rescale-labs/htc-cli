package tabler

import (
	"time"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
)

// Generic function for handling oapi.OptInstant, OptDateTime, and
// OptNilDateTime.
//
// I'm honestly amazed this worked at all. But, I guess it makes enough
// sense to keep it.

// First name the union so we can use it in a few places.
type anyTime interface{ time.Time | oapi.Instant }

// Then construct the generic interface as a type constraint.
type TimeGet[T anyTime] interface {
	Get() (T, bool)
}

// Finally, implement the function and have it work for either
// time type. We succeed b/c oapi.Instant has the same base type as
// time.Time.
func formatDateTime[T anyTime](i TimeGet[T]) string {
	t, ok := i.Get()
	if !ok {
		return ""
	}
	return time.Time(t).Format(time.DateTime)
}

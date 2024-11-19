package tabler

import (
	"fmt"
	"io"

	oapi "github.com/rescale/htc-storage-cli/v2/api/_oas"
)

type HTCJobs []oapi.HTCJob

func (s HTCJobs) Fields() []Field {
	return []Field{
		Field{"UUID", "%-38s", "%-38s"},
		Field{"Created", "%19s", "%19s"},
		Field{"Completed", "%19s", "%19s"},
		Field{"Status", "%12s", "%12s"},
	}
}

func (s HTCJobs) WriteRows(rowFmt string, w io.Writer) error {
	for _, j := range s {
		_, err := fmt.Fprintf(
			w, rowFmt,
			j.JobUUID.Value,
			formatDateTime(j.CreatedAt),
			formatDateTime(j.CompletedAt),
			j.Status.Value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

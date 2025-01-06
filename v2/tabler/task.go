package tabler

import (
	"fmt"
	"io"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
)

type HTCTasks []oapi.HTCTask

func (s HTCTasks) Fields() []Field {
	return []Field{
		Field{"ID", "%-38s", "%-38s"},
		Field{"Name", "%24s", "%24.24s"},
		Field{"Created", "%19s", "%19s"},
		Field{"Last Active", "%19s", "%19s"},
		Field{"Archived", "%19s", "%19s"},
	}
}

func (s HTCTasks) WriteRows(rowFmt string, w io.Writer) error {
	for _, t := range s {
		_, err := fmt.Fprintf(
			w, rowFmt,
			t.TaskId.Value,
			t.TaskName,
			formatDateTime(t.CreatedAt),
			formatDateTime(t.LastActiveAt),
			formatDateTime(t.ArchivedAt),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

package tabler

import (
	"fmt"
	"io"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
)

type HTCJobLog []oapi.HTCLogEvent

func (s HTCJobLog) Fields() []Field {
	return []Field{
		Field{"Timestamp", "%-38s", "%-38s"},
		Field{"Message", "%19s", "%19s"},
	}
}

func (s HTCJobLog) WriteRows(rowFmt string, w io.Writer) error {
	for _, t := range s {
		_, err := fmt.Fprintf(
			w, rowFmt,
			formatDateTime(t.Timestamp),
			t.Message.Value,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

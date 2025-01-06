package tabler

import (
	"fmt"
	"io"

	oapi "github.com/rescale-labs/htc-cli/v2/api/_oas"
)

type HTCProjects []oapi.HTCProject

func (s HTCProjects) Fields() []Field {
	return []Field{
		Field{"ID", "%-36s", "%-36s"},
		Field{"Name", "%-24s", "%-24.24s"},
		Field{"Description", "%-24s", "%-24.24s"},
		Field{"Created", "%19s", "%19s"},
	}
}

func (s HTCProjects) WriteRows(rowFmt string, w io.Writer) error {
	for _, p := range s {
		_, err := fmt.Fprintf(
			w, rowFmt,
			p.ProjectId.Value,
			p.ProjectName,
			p.ProjectDescription,
			formatDateTime(p.CreatedAt),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

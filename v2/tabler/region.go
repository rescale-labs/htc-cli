package tabler

import (
	"fmt"
	"io"
)

var regionFields = []Field{
	Field{"Region", "%-38s", "%-38s"},
}

type Regions []string

func (r Regions) Fields() []Field {
	return regionFields
}

func (r Regions) WriteRows(rowFmt string, w io.Writer) error {
	for _, s := range r {
		if _, err := fmt.Fprintf(w, rowFmt, s); err != nil {
			return err
		}
	}
	return nil
}

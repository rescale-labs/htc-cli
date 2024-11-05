// A simple package for writing API response data as text tables to
// stdout.
package tabler

import (
	"fmt"
	"io"
	"strings"
)

type Field struct {
	Name         string
	HeaderFormat string
	CellFormat   string
}

type Tabler interface {
	// Returns all fields to print in the table, in order.
	Fields() []Field

	// Takes a row printf format string and a writer. Writes all rows to
	// the writer, returning the first error encountered, if any.
	WriteRows(string, io.Writer) error
}

func WriteTable(t Tabler, w io.Writer) error {
	fields := t.Fields()
	var headerFmts, cellFmts []string
	var headers []any
	for _, f := range fields {
		headerFmts = append(headerFmts, f.HeaderFormat)
		cellFmts = append(cellFmts, f.CellFormat)
		headers = append(headers, strings.ToUpper(f.Name))
	}
	headerFmt := strings.Join(headerFmts, "  ") + "\n"
	rowFmt := strings.Join(cellFmts, "  ") + "\n"

	if _, err := fmt.Fprintf(w, headerFmt, headers...); err != nil {
		return err
	}
	return t.WriteRows(rowFmt, w)
}

// A simple package for writing API response data as text tables to
// stdout.
package tabler

import (
	"fmt"
	"io"
	"log"
	"strings"
)

type Field struct {
	Name         string
	HeaderFormat string
	CellFormat   string
}

type Tabler interface {
	Fields() []Field
	Rows() [][]any
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
	cellFmt := strings.Join(cellFmts, "  ") + "\n"
	log.Printf("headerFmt=%q cellFmt=%q", headerFmt, cellFmt)
	if _, err := fmt.Fprintf(w, headerFmt, headers...); err != nil {
		return err
	}

	for _, values := range t.Rows() {
		if _, err := fmt.Fprintf(w, cellFmt, values...); err != nil {
			return err
		}
	}

	return nil
}

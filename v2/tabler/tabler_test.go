package tabler

import (
	"bytes"
	"math"
	"strings"
	"testing"
)

type ThreeColTabler struct{}

func (t *ThreeColTabler) Fields() []Field {
	return []Field{
		Field{"name", "%-10s", "%-10s"},
		Field{"length", "%12s", "%12d"},
		Field{"float", "%6s", "%6.2f"},
	}
}

func (t ThreeColTabler) Rows() [][]any {
	return [][]any{
		[]any{"Bruce", 12919, 2.3443},
		[]any{"Jim", 35, 322.64},
		[]any{"John", 393, math.Pi},
	}
}

func TestWriteTable(t *testing.T) {
	var w bytes.Buffer

	tabler := ThreeColTabler{}
	err := WriteTable(&tabler, &w)
	if err != nil {
		t.Errorf("WriteTable: wanted nil, got %s", err)
	}

	expected := strings.TrimSpace(`
NAME              LENGTH   FLOAT
Bruce              12919    2.34
Jim                   35  322.64
John                 393    3.14
`) + "\n"

	actual := w.String()
	t.Logf("Output:\n%s", actual)
	if expected != actual {
		t.Errorf("Wanted output:\n%s\n\nGot:\n\n%s", expected, actual)
	}
}

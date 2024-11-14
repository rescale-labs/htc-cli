package tabler

import (
	"fmt"
	"io"

	"github.com/rescale/htc-storage-cli/v2/config"
)

type ContextConf struct {
	Name     string `json:"name"`
	Selected bool   `json:"selected"`
	*config.ContextConf
}

type ContextConfs []*ContextConf

func (c ContextConfs) Fields() []Field {
	return []Field{
		Field{"", "%-2s", "%2s"},
		Field{"Name", "%-38s", "%-38s"},
		Field{"Project ID", "%19s", "%19s"},
		Field{"Task ID", "%19s", "%19s"},
	}
}

func (c ContextConfs) WriteRows(rowFmt string, w io.Writer) error {
	for _, i := range c {
		var selected string
		if i.Selected {
			selected = "*"
		}
		_, err := fmt.Fprintf(
			w, rowFmt,
			selected,
			i.Name,
			i.ProjectId,
			i.TaskId,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

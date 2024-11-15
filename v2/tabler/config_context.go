package tabler

import (
	"fmt"
	"io"

	"github.com/rescale/htc-storage-cli/v2/config"
)

type ContextIdentity struct {
	Email         string `json:"email"`
	WorkspaceId   string `json:"workspaceId"`
	WorkspaceName string `json:"workspaceName"`
}

type ContextConf struct {
	Name     string          `json:"name"`
	Selected bool            `json:"selected"`
	Identity config.Identity `json:"identity,omitempty"`
	*config.ContextConf
}

type ContextConfs []*ContextConf

func (c ContextConfs) Fields() []Field {
	var nameLen, wsLen, emailLen int
	for _, i := range c {
		nameLen = max(nameLen, len(i.Name))
		wsLen = max(wsLen, len(i.Identity.WorkspaceName))
		emailLen = max(emailLen, len(i.Identity.Email))
	}
	nameField := fmt.Sprintf("%%-%ds", nameLen)
	wsField := fmt.Sprintf("%%-%ds", wsLen)
	emailField := fmt.Sprintf("%%-%ds", emailLen)
	return []Field{
		Field{"", "%-2s", "%2s"},
		Field{"Name", nameField, nameField},
		Field{"Workspace ID", "%12s", "%12s"},
		Field{"Workspace Name", wsField, wsField},
		Field{"Email", emailField, emailField},
		Field{"Project ID", "%-36s", "%-36s"},
		Field{"Task ID", "%-36s", "%-36s"},
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
			i.Identity.WorkspaceId,
			i.Identity.WorkspaceName,
			i.Identity.Email,
			i.ProjectId,
			i.TaskId,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

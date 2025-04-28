package workspace

import (
	"github.com/spf13/cobra"
)

var WorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Commands for managing HTC workspaces",
}

func init() {
	WorkspaceCmd.AddCommand(ClustersCmd)
	WorkspaceCmd.AddCommand(RetentionPolicyCmd)
}

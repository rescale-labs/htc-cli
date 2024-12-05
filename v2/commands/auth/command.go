package auth

import (
	"github.com/spf13/cobra"
)

var AuthCmd = &cobra.Command{
	Use: "auth",
}

func init() {
	AuthCmd.AddCommand(LoginCmd)
	AuthCmd.AddCommand(LogoutCmd)
	AuthCmd.AddCommand(WhoAmICmd)
}

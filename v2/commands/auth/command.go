package auth

import (
	"github.com/spf13/cobra"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Log in and out of HTC",
}

func init() {
	AuthCmd.AddCommand(LoginCmd)
	AuthCmd.AddCommand(LogoutCmd)
	AuthCmd.AddCommand(WhoAmICmd)
}

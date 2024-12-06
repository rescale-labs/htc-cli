package auth

import (
	"github.com/rescale/htc-storage-cli/v2/common"
	"github.com/spf13/cobra"
)

func Logout(cmd *cobra.Command, args []string) error {
	r, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}
	return r.Config.DeleteCredentials()
}

var LogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Removes any stored credentials or token for this context.",
	// Long:
	Run: common.WrapRunE(Logout),
}

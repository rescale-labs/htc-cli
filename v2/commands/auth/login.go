package auth

import (
	"github.com/rescale/htc-storage-cli/v2/common"
	"github.com/spf13/cobra"
)

func Login(cmd *cobra.Command, args []string) error {
	r, err := common.NewRunner(cmd)
	if err != nil {
		return err
	}

	if err := r.RenewToken(); err != nil {
		return err
	}
	return nil
}

var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Obtains a bearer token (JWT) from the HTC API using RESCALE_API_KEY",
	// Long:
	Run: common.WrapRunE(Login),
}

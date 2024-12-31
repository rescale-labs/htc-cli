// Simple command for printing our tagged version.
//
// NB: much more is available in debug/buildinfo or
// by running `go version -m ./build/htc`. But, even though
// the go build tool icnludes the git SHA, it doesn't
// include the tagged version. So...ldflags it remains.
package version

import (
	"fmt"
	"runtime"

	"github.com/rescale/htc-storage-cli/v2/common"
	"github.com/spf13/cobra"
)

// Set by go build -ldflags="-X ..." in Makefile for all distributed
// builds.
var Version = "devel"

var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version of this CLI tool",
	Run: common.WrapRunE(
		func(cmd *cobra.Command, _ []string) error {
			_, err := fmt.Printf(
				"htc CLI version: %s (%v %s/%s)\n",
				Version,
				runtime.Version(), runtime.GOOS, runtime.GOARCH)
			return err
		},
	),
}

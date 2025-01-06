package metrics

import (
	"github.com/spf13/cobra"
)

var MetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Fetch HTC metrics in OpenMetrics format",
}

func init() {
	MetricsCmd.AddCommand(GetCmd)
}

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zigelboim-misha/ci-templates-metrics/receiver"
)

// receiverCmd represents the receiver command
var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "Running the receiver µservice which will wait for HTTP requests and then exports its metrics in a prometheus format",
	Long: `Running the receiver µservice.
When a new CI job name arrives, it will first register a new prometheus metric and then increase it.
If an existing CI job arrives, it will only increment what was requested.`,
	Run: runReceiver,
}

func runReceiver(cmd *cobra.Command, args []string) {
	receiver.RunReceiver()
}

func init() {
	rootCmd.AddCommand(receiverCmd)
}

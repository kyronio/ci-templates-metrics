package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zigelboim-misha/ci-templates-metrics/exporter"
)

var (
	filePath     string
	receiverHost string
)

// exporterCmd represents the exporter command
var exporterCmd = &cobra.Command{
	Use:   "exporter [file-path] [receiver-hostname]",
	Short: "Running the metrics exporter µservice that will send the received metrics to the given host",
	Long: `The exporter command will wait for a .json file to be written by your CI job.
When the files is created, the µservice will send it to the given receiver µservice hostname.`,
	Run:  runExporter,
	Args: cobra.ExactArgs(2),
}

func runExporter(cmd *cobra.Command, args []string) {
	filePath = args[0]
	receiverHost = args[1]

	exporter.RunExporter(filePath, receiverHost)
}

func init() {
	rootCmd.AddCommand(exporterCmd)
}

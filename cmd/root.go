package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "1.0.0-alpha"

var rootCmd = &cobra.Command{
	Use:     "sigil",
	Short:   "Certificate stamping and dispatch system for Amrita events",
	Version: Version,
	Long: `Sigil is a certificate dispatch system for Amrita events.
It handles three types of certificates:
  - Certificate of Participation (external)
  - Certificate of Achievement (external)
  - Certificate of Recognition (internal)

The system uses RabbitMQ for reliable queue-based processing and
pdfcpu for PDF certificate stamping.

Run 'sigil --help' for available commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running sigil version: %s\n", Version)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

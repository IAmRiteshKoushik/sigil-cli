package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "1.0.0-alpha"

var rootCmd = &cobra.Command{
	Use:     "sigil-cli",
	Short:   "The CLI for dispatching certificates",
	Version: Version,
	Long: `Sigil CLI is a simple tool to orchestrate sending certificates in 
Amrita events. It relies on github.com/pdfcpu/pdfcpu to do the heavy lifting 
of certificate generatoin and uses RabbitMQ for reliability during dispatches.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running sigil-cli version: %s\n", Version)
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

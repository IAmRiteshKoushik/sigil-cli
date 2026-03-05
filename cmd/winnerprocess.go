package cmd

import (
	"fmt"
	"log"

	"github.com/IAmRiteshKoushik/sigil/pkg"
	"github.com/spf13/cobra"
)

var winnerprocessCmd = &cobra.Command{
	Use:   "winnerprocess",
	Short: "Process CSV containing data for winners",
	Args:   cobra.ExactArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		csvFile := args[0]

		winners, err := pkg.ParseWinnersCSVFile(csvFile)
		if err != nil {
			log.Fatalf("Error parsing CSV file: %v", err)
		}
		if len(winners) == 0 {
			fmt.Println("No valid winners records found")
			return
		}

		queueName := "cert_winners"
		if err := pkg.PublishToWinnersQueue(queueName, winners); err != nil {
			log.Fatalf("Error publishing to queue: %v", err)
		}
		fmt.Println("Winners CSV processing complete")
	},
}

func init() {
	rootCmd.AddCommand(winnerprocessCmd)
}

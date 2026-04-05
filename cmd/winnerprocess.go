package cmd

import (
	"fmt"
	"log"

	"github.com/IAmRiteshKoushik/sigil/pkg"
	"github.com/spf13/cobra"
)

var winnerprocessCmd = &cobra.Command{
	Use:   "winnerprocess [csv-file]",
	Short: "Process winners CSV and publish to cert_winners queue",
	Long: `Reads winner data from a CSV file and publishes to the cert_winners queue.
Each winner record includes: event_name, position, student_name, student_email.
This command is the first step in the Certificate of Achievement workflow.

CSV format:
  event_name,position,student_name,student_email
  AGAMOTTO PROTOCOL,1st Place,John Doe,john@example.com

Example:
  sigil winnerprocess storage/winners.csv`,
	Args: cobra.ExactArgs(1),
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

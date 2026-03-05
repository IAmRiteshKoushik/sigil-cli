package cmd

import (
	"fmt"
	"log"

	"github.com/IAmRiteshKoushik/sigil/pkg"
	"github.com/spf13/cobra"
)

var partprocessCmd = &cobra.Command{
	Use:   "partprocess [csv-file]",
	Short: "Process CSV file and add student data to participation certificate queue",
	Long:  `Read student data from CSV file and publish to cert_ queue as JSON payloads`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csvFile := args[0]

		students, err := pkg.ParseCSVFile(csvFile)
		if err != nil {
			log.Fatalf("Error parsing CSV file: %v", err)
		}

		if len(students) == 0 {
			fmt.Println("No valid student records found")
			return
		}

		eventName := pkg.ExtractEventName(csvFile)
		queueName := fmt.Sprintf("cert_%s", eventName)

		fmt.Printf("Processing %d students for event: %s\n", len(students), eventName)
		fmt.Printf("Publishing to queue: %s\n", queueName)

		if err := pkg.PublishToQueue(queueName, students); err != nil {
			log.Fatalf("Error publishing to queue: %v", err)
		}

		fmt.Println("CSV processing completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(partprocessCmd)
}

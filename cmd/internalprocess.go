package cmd

import (
	"fmt"
	"log"

	"github.com/IAmRiteshKoushik/sigil/pkg"
	"github.com/spf13/cobra"
)

var internalprocessCmd = &cobra.Command{
	Use:   "internalprocess [csv-file]",
	Short: "Process CSV file and add internal student data to recognition queue",
	Long:  `Read student data from CSV file and publish to cert_ queue as JSON payloads`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csvFile := args[0]

		students, err := pkg.ParseInternalCSVFile(csvFile)
		if err != nil {
			log.Fatalf("Error parsing CSV file: %v", err)
		}

		if len(students) == 0 {
			fmt.Println("No valid student records found")
			return
		}

		queueName := "cert_internal"
		if err := pkg.PublishToQueue(queueName, students); err != nil {
			log.Fatalf("Error publishing to queue: %v", err)
		}
		fmt.Println("Internal CSV processing complete")
	},
}

func init() {
	rootCmd.AddCommand(internalprocessCmd)
}

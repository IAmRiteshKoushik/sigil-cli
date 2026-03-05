package cmd

import (
	"fmt"
	"log"

	"github.com/IAmRiteshKoushik/sigil/pkg"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [events-file]",
	Short: "Create RabbitMQ queues from events.txt",
	Long:  `Read events from a file and create cert_ and dispatch_ queues for each event`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		eventsFile := args[0]

		events, err := pkg.ReadEventsFile(eventsFile)
		if err != nil {
			log.Fatalf("Error reading events file: %v", err)
		}

		fmt.Printf("Found %d events\n", len(events))

		if err := pkg.CreateQueues(events); err != nil {
			log.Fatalf("Error creating queues: %v", err)
		}

		fmt.Println("Queue creation completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

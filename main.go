package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	LoadConfig()
	setupLogger()
}

var createCmd = &cobra.Command{
	Use:   "create [events-file]",
	Short: "Create RabbitMQ queues for events",
	Long:  `Read events from a file and create cert_ and dispatch_ queues for each event`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		eventsFile := args[0]

		events, err := readEventsFile(eventsFile)
		if err != nil {
			log.Fatalf("Error reading events file: %v", err)
		}

		fmt.Printf("Found %d events\n", len(events))

		if err := createQueues(events); err != nil {
			log.Fatalf("Error creating queues: %v", err)
		}

		fmt.Println("Queue creation completed successfully")
	},
}

var processCmd = &cobra.Command{
	Use:   "process [csv-file]",
	Short: "Process CSV file and add student data to certificate queue",
	Long:  `Read student data from CSV file and publish to cert_ queue as JSON payloads`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		csvFile := args[0]

		students, err := parseCSVFile(csvFile)
		if err != nil {
			log.Fatalf("Error parsing CSV file: %v", err)
		}

		if len(students) == 0 {
			fmt.Println("No valid student records found")
			return
		}

		eventName := extractEventName(csvFile)
		queueName := fmt.Sprintf("cert_%s", eventName)

		fmt.Printf("Processing %d students for event: %s\n", len(students), eventName)
		fmt.Printf("Publishing to queue: %s\n", queueName)

		if err := publishToQueue(queueName, students); err != nil {
			log.Fatalf("Error publishing to queue: %v", err)
		}

		fmt.Println("CSV processing completed successfully")
	},
}

var processBatchCmd = &cobra.Command{
	Use:   "process-batch [reports-folder]",
	Short: "Process all CSV files in a folder",
	Long:  `Read all CSV files from the reports folder and process them in bulk`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		reportsDir := args[0]
		moveProcessed, _ := cmd.Flags().GetBool("move")

		// Check if reports directory exists
		if _, err := os.Stat(reportsDir); os.IsNotExist(err) {
			log.Fatalf("Reports directory '%s' does not exist", reportsDir)
		}

		// Find all CSV files
		csvFiles, err := findCSVFiles(reportsDir)
		if err != nil {
			log.Fatalf("Error finding CSV files: %v", err)
		}

		if len(csvFiles) == 0 {
			fmt.Printf("No CSV files found in '%s'\n", reportsDir)
			return
		}

		fmt.Printf("Found %d CSV files to process\n", len(csvFiles))
		fmt.Println("Starting batch processing...")

		processed := 0
		failed := 0

		// Process each CSV file
		for i, csvFile := range csvFiles {
			filename := filepath.Base(csvFile)
			fmt.Printf("\n[%d/%d] Processing: %s\n", i+1, len(csvFiles), filename)

			students, err := parseCSVFile(csvFile)
			if err != nil {
				fmt.Printf("❌ Error parsing %s: %v\n", filename, err)
				failed++
				continue
			}

			if len(students) == 0 {
				fmt.Printf("⚠️  No valid student records found in %s\n", filename)
				failed++
				continue
			}

			eventName := extractEventName(csvFile)
			queueName := fmt.Sprintf("cert_%s", eventName)

			fmt.Printf("📊 Processing %d students for event: %s\n", len(students), eventName)
			fmt.Printf("📤 Publishing to queue: %s\n", queueName)

			if err := publishToQueue(queueName, students); err != nil {
				fmt.Printf("❌ Error publishing %s to queue: %v\n", filename, err)
				failed++
				continue
			}

			fmt.Printf("✅ Successfully processed: %s\n", filename)
			processed++

			// Move processed file if flag is set
			if moveProcessed {
				processedDir := filepath.Join(reportsDir, "processed")
				if err := os.MkdirAll(processedDir, 0755); err != nil {
					fmt.Printf("⚠️  Warning: Could not create processed directory: %v\n", err)
				} else {
					newPath := filepath.Join(processedDir, filename)
					if err := os.Rename(csvFile, newPath); err != nil {
						fmt.Printf("⚠️  Warning: Could not move file %s: %v\n", filename, err)
					} else {
						fmt.Printf("📁 Moved processed file to: processed/%s\n", filename)
					}
				}
			}
		}

		// Summary
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Printf("Batch Processing Complete\n")
		fmt.Printf("✅ Successfully processed: %d files\n", processed)
		fmt.Printf("❌ Failed: %d files\n", failed)

		if failed > 0 {
			fmt.Printf("\n⚠️  Warning: %d files failed to process\n", failed)
			os.Exit(1)
		} else {
			fmt.Printf("\n🎉 All files processed successfully!\n")
		}
	},
}

func findCSVFiles(dir string) ([]string, error) {
	var csvFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".csv" {
			csvFiles = append(csvFiles, path)
		}

		return nil
	})

	return csvFiles, err
}

func init() {
	processBatchCmd.Flags().Bool("move", false, "Move processed files to 'processed' subfolder")
}

func main() {
	var rootCmd = &cobra.Command{Use: "sigil"}
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(processBatchCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

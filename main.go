package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type Config struct {
	RabbitMQ struct {
		URL string `toml:"url"`
	} `toml:"rabbitmq"`
	CLI struct {
		LogLevel string `toml:"log_level"`
	} `toml:"cli"`
}

var cfg Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
}

func createQueues(events []string) error {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	for _, event := range events {
		if strings.TrimSpace(event) == "" {
			continue
		}

		queues := []string{
			fmt.Sprintf("cert_%s", event),
			fmt.Sprintf("dispatch_%s", event),
		}

		for _, queueName := range queues {
			_, err := ch.QueueDeclare(
				queueName, // name
				true,      // durable
				false,     // delete when unused
				false,     // exclusive
				false,     // no-wait
				nil,       // arguments
			)
			if err != nil {
				log.Printf("Failed to declare queue %s: %v", queueName, err)
				continue
			}
			fmt.Printf("Created queue: %s\n", queueName)
		}
	}

	return nil
}

func readEventsFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	var events []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			events = append(events, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %v", filename, err)
	}

	return events, nil
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

func main() {
	var rootCmd = &cobra.Command{Use: "sigil"}
	rootCmd.AddCommand(createCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

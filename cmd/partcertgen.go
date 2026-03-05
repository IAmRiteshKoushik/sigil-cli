package cmd

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
)

type Event struct {
	StudentName  string `json:"student_name"`
	StudentEmail string `json:"student_email"`
	EventName    string `json:"event_name"`
}

var partcertgenCmd = &cobra.Command{
	Use:   "partcertgen [queue_base_name] [template_path]",
	Args:  cobra.ExactArgs(2),
	Short: "Generates certificates based on events from a RabbitMQ queue",
	Long: `The certgen command connects to a RabbitMQ queue, dequeues an event
containing student information, creates a certificate file, and dispatches
a new event to a dispatch queue.

Usage:
  certgen my_event_queue

This will connect to "cert_my_event_queue", process an event, create a file
in @storage/cert/event_name/student_email, and publish to "dispatch_event_name".`,
	Run: func(cmd *cobra.Command, args []string) {
		queueBaseName := args[0]
		templatePath := args[1]

		certQueueName := "cert_" + queueBaseName
		dispatchQueuePrefix := "dispatch_"

		storagePath := "/home/rk/ws/dev/sigil/storage/cert/"

		conn, err := amqp.Dial("amqp://admin:admin123@localhost:5672/sigil-vhost")
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %v", err)
		}
		defer ch.Close()

		// Create directory if it doesn't exist
		eventDirPath := filepath.Join(storagePath, queueBaseName)
		err = os.MkdirAll(eventDirPath, 0755)
		if err != nil {
			log.Printf("Error creating directory %s: %v", eventDirPath, err)
			return
		}

		q, err := ch.QueueDeclare(
			certQueueName, // name
			true,          // durable
			false,         // delete when unused
			false,         // exclusive
			false,         // no-wait
			nil,           // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %v", err)
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Fatalf("Failed to register a consumer: %v", err)
		}

		var forever chan struct{}

		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)

				var event Event
				err := json.Unmarshal(d.Body, &event)
				if err != nil {
					log.Printf("Error unmarshalling JSON: %v", err)
					continue
				}

				// Process the event name
				eventNameClean := strings.ReplaceAll(event.EventName, "_", " ")
				eventName := strings.ToUpper(eventNameClean)
				studentName := strings.ToUpper(event.StudentName)
				outputPath := eventDirPath + "/" + event.StudentEmail + ".pdf"

				// Create and write to file
				// cmd := fmt.Sprintf("just stamp-participation \"%s\" \"%s\" \"%s\" \"%s\"", studentName, eventName, templatePath, outputPath)
				cmd := exec.Command("just",
					"stamp-participation",
					studentName,
					eventName,
					templatePath,
					outputPath,
				)
				output, err := cmd.CombinedOutput()
				if err != nil {
					log.Fatalf("Command failed: %s\nOutput: %s", err, string(output))
				}

				// dispatchQueueName := dispatchQueuePrefix + event.EventName
				dispatchQueueName := dispatchQueuePrefix + "internal"
				err = ch.Publish(
					"",                // exchange
					dispatchQueueName, // routing key
					false,             // mandatory
					false,             // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        d.Body,
					})
				if err != nil {
					log.Printf("Failed to publish a message to dispatch queue: %v", err)
				} else {
					log.Printf("Event dispatched to queue: %s", dispatchQueueName)
				}
			}
		}()

		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	},
}

func init() {
	rootCmd.AddCommand(partcertgenCmd)
}

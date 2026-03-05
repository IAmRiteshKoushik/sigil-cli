package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"

	pkglib "github.com/IAmRiteshKoushik/sigil/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
	"gopkg.in/gomail.v2"
)

var internalcertsendCmd = &cobra.Command{
	Use:   "internalcertsend <event_name>",
	Short: "Dequeues events from RabbitMQ and sends recognition certificates via email.",
	Long: `This command dequeues events from a RabbitMQ queue named "dispatch_internal".
For each event, it constructs an email with a participation certificate attached.
The certificate file is expected to be found at @storage/cert/internal/<student_email>.pdf.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		eventName := args[0]
		queueName := "dispatch_" + eventName

		pkglib.LoadConfig()
		config := pkglib.GetConfig()

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

		q, err := ch.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %v", err)
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Fatalf("Failed to register a consumer: %v", err)
		}

		log.Printf(" [*] Waiting for messages in queue %s. To exit press CTRL+C", queueName)

		for d := range msgs {
			var event Event
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				d.Nack(false, false)
				continue
			}

			log.Printf("Received event for student %s (%s) for event %s", event.StudentName, event.StudentEmail, event.EventName)

			certPath := filepath.Join(config.StorageDir, "cert", "internal", event.StudentEmail+".pdf")

			subject := "Anokha 2026 - Certificate of Recognition"

			body := fmt.Sprintf("Hey %s,\n\nCongratulations on successfully organizing Anokha 2026! On behalf of the organizing committee at Amrita Vishwa Vidyapeetham, we want to extend our heartfelt thanks to you for being part of this year's journey. Your energy, innovation, and passion were what made this edition of Anokha truly vibrant and memorable. \n\nAttached to this email, you will find your Certificate of Recognition. This document serves as a formal recognition of the hard work and skill you brought to the event.\n\nBest Regards\n\nThe Anokha 2026 - Organizing Committee\nAmrita Vishwa Vidyapeetham", event.StudentName)

			for {
				message := gomail.NewMessage()
				message.SetHeader("From", config.SMTPUser)
				message.SetHeader("To", event.StudentEmail)
				message.SetHeader("Subject", subject)

				message.SetBody("text/plain", body)
				message.Attach(certPath, gomail.Rename("recognition-certificate.pdf"))

				dailer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPass)

				if err := dailer.DialAndSend(message); err != nil {
					log.Printf("Error: %v", err)
					panic(err)
				}

				break
			}

			log.Printf("Successfully sent certificate to %s", event.StudentEmail)
			d.Ack(false)
		}
	},
}

func init() {
	rootCmd.AddCommand(internalcertsendCmd)
}

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

var winnersendCmd = &cobra.Command{
	Use:   "winnersend <name>",
	Short: "Dequeue and send Certificate of Achievement via email",
	Long: `Consumes messages from the dispatch_<name> queue and sends
Certificate of Achievement emails to winners.

The certificate PDF is expected at: storage/cert/winners/<student_email>.pdf

Each email includes:
  - Subject: "Anokha 2026 - Certificate of Achievement"
  - Personalized congratulations body
  - Certificate PDF attachment renamed to "certification-of-achievement.pdf"

Example:
  sigil winnersend winners`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		queueName := "dispatch_" + name

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
			var event pkglib.WinnerData
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				d.Nack(false, false)
				continue
			}

			log.Printf("Received event for student %s (%s) for event %s", event.StudentName, event.StudentEmail, event.EventName)

			certPath := filepath.Join(config.StorageDir, "cert", "winners", event.StudentEmail+".pdf")
			subject := "Anokha 2026 - Certificate of Achievement"

			body := fmt.Sprintf("Hey %s,\n\nCongratulations on your outstanding performance at Anokha 2026! We are thrilled to recognize your exceptional skill and dedication, which sets a high standard for innovation within the Amrita Vishwa Vidyapeetham community.\n\nIn recognition of your significant accomplishment, we have attached your Certificate of Achievement to this email. This honor celebrates your success and the excellence you demonstrated during the event.\n\nWe look forward to seeing your future breakthroughs.\n\nBest Regards,\n\nThe Anokha 2026 - Organizing Committee\nAmrita Vishwa Vidyapeetham", event.StudentName)

			for {
				message := gomail.NewMessage()
				message.SetHeader("From", config.SMTPUser)
				message.SetHeader("To", event.StudentEmail)
				message.SetHeader("Subject", subject)

				message.SetBody("text/plain", body)
				message.Attach(certPath, gomail.Rename("certification-of-achievement.pdf"))

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
	rootCmd.AddCommand(winnersendCmd)
}

package pkg

import (
	"fmt"
	"log"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CreateQueues(events []string) error {
	cfg := GetConfig()

	conn, err := amqp.Dial(cfg.RabbitMQURL)
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
			// Check if queue already exists
			exists, err := QueueExists(ch, queueName)
			if err != nil {
				log.Printf("Failed to check if queue %s exists: %v", queueName, err)
				continue
			}

			if exists {
				log.Printf("Queue %s already exists, skipping creation...", queueName)
				fmt.Printf("Queue %s already exists, skipping...\n", queueName)
				continue
			}

			// Create queue if it doesn't exist
			_, err = ch.QueueDeclare(
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
			log.Printf("Created queue: %s\n", queueName)
		}
	}

	return nil
}

func QueueExists(ch *amqp.Channel, queueName string) (bool, error) {
	// Try to passively declare the queue to check if it exists
	_, err := ch.QueueDeclarePassive(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		if strings.Contains(err.Error(), "NOT_FOUND") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

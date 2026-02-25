package pkg

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type StudentData struct {
	StudentName  string `json:"student_name"`
	StudentEmail string `json:"student_email"`
	EventName    string `json:"event_name"`
}

func ParseCSVFile(filename string) ([]StudentData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file %s: %v", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file %s: %v", filename, err)
	}

	var students []StudentData
	for i, record := range records {
		// Skip header if it exists (check if first row looks like headers)
		if i == 0 && (strings.ToLower(strings.TrimSpace(record[0])) == "student_name" ||
			strings.ToLower(strings.TrimSpace(record[1])) == "student_email") {
			continue
		}

		// Ensure we have at least 2 columns
		if len(record) < 2 {
			log.Printf("Skipping row %d: insufficient columns", i+1)
			continue
		}

		student := StudentData{
			StudentName:  strings.TrimSpace(record[0]),
			StudentEmail: strings.TrimSpace(record[1]),
			EventName:    ExtractEventName(filename),
		}

		if student.StudentName == "" || student.StudentEmail == "" {
			log.Printf("Skipping row %d: empty name or email", i+1)
			continue
		}

		students = append(students, student)
	}

	return students, nil
}

func ExtractEventName(filename string) string {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func PublishToQueue(queueName string, students []StudentData) error {
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

	for _, student := range students {
		payload := fmt.Sprintf(`{"student_name":"%s","student_email":"%s","event_name":"%s"}`,
			student.StudentName, student.StudentEmail, student.EventName)

		err := ch.Publish(
			"",        // exchange
			queueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(payload),
			})
		if err != nil {
			log.Printf("Failed to publish student %s to queue %s: %v",
				student.StudentName, queueName, err)
			continue
		}

		log.Printf("Published student %s to queue %s", student.StudentName, queueName)
	}

	return nil
}

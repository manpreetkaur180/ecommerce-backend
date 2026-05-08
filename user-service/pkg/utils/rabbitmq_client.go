package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const UserRegisteredQueue = "user.registered"
const OTPEmailQueue = "email.otp"
const VerificationEmailQueue = "email.verification"

func PublishUserRegistered(name, email, phone string) error {
	payload := map[string]string{
		"name":  name,
		"email": email,
		"phone": phone,
	}

	return publishWithRetry(UserRegisteredQueue, payload)
}

func PublishOTPEmail(name, email, otp string) error {
	payload := map[string]string{
		"name":  name,
		"email": email,
		"otp":   otp,
	}

	return publishWithRetry(OTPEmailQueue, payload)
}

func PublishVerificationEmail(name, email, link string) error {
	payload := map[string]string{
		"name":  name,
		"email": email,
		"link":  link,
	}

	return publishWithRetry(VerificationEmailQueue, payload)
}

func publishWithRetry(queueName string, payload map[string]string) error {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	var lastErr error
	for attempt := 1; attempt <= 12; attempt++ {
		if err := publishMessage(rabbitURL, queueName, payload); err != nil {
			lastErr = err
			time.Sleep(5 * time.Second)
			continue
		}

		return nil
	}

	return fmt.Errorf("failed to publish %s after retries: %w", queueName, lastErr)
}

func publishMessage(rabbitURL, queueName string, payload map[string]string) error {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

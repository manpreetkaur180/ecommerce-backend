package consumers

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"notification-service/internal/models"
	"notification-service/internal/services"

	amqp "github.com/rabbitmq/amqp091-go"
)

const UserRegisteredQueue = "user.registered"
const OTPEmailQueue = "email.otp"
const VerificationEmailQueue = "email.verification"

func StartEmailConsumers(service *services.NotificationService) {
	go startConsumer(UserRegisteredQueue, func(body []byte) error {
		var event models.UserRegisteredEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		if err := service.HandleUserRegistered(event.Name, event.Email); err != nil {
			return err
		}

		log.Println("Welcome email processed for", event.Email)
		return nil
	})

	go startConsumer(OTPEmailQueue, func(body []byte) error {
		var event models.OTPEmailEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		if err := service.HandleOTPEmail(event.Name, event.Email, event.OTP); err != nil {
			return err
		}

		log.Println("OTP email processed for", event.Email)
		return nil
	})

	go startConsumer(VerificationEmailQueue, func(body []byte) error {
		var event models.VerificationEmailEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		if err := service.HandleVerificationEmail(event.Name, event.Email, event.Link); err != nil {
			return err
		}

		log.Println("Verification email processed for", event.Email)
		return nil
	})
}

func startConsumer(queueName string, handler func([]byte) error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	for {
		if err := consumeQueue(rabbitURL, queueName, handler); err != nil {
			log.Println("RabbitMQ consumer error for", queueName+":", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func consumeQueue(rabbitURL, queueName string, handler func([]byte) error) error {
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

	messages, err := ch.Consume(
		queue.Name,
		"notification-service-"+queue.Name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("RabbitMQ consumer started for queue", queue.Name)

	for message := range messages {
		if err := handler(message.Body); err != nil {
			log.Println("Failed to process", queue.Name, "event:", err)
			message.Nack(false, false)
			continue
		}

		message.Ack(false)
	}

	return nil
}

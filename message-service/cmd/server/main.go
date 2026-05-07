package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"message-service/config"
	"message-service/internal/handlers"
	smtpprovider "message-service/internal/providers/smtp"
	"message-service/internal/services"
)

func main() {

	app := fiber.New()

	cfg := config.LoadSMTPConfig()

	smtpProvider := smtpprovider.NewSMTPProvider(cfg)

	messageService := services.NewMessageService(
		smtpProvider,
	)

	emailHandler := handlers.NewEmailHandler(
	messageService,
)

app.Post("/email/send", emailHandler.SendEmail)

	port := os.Getenv("PORT")

	if port == "" {
		port = "3002"
	}

	log.Println("Message service running on", port)

	log.Fatal(app.Listen(":" + port))
}
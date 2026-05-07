package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"message-service/config"
	"message-service/internal/handlers"

	smtpprovider "message-service/internal/providers/smtp"
	twilioprovider "message-service/internal/providers/twilio"

	"message-service/internal/services"
)

func main() {

	app := fiber.New()

	// ---------------- SMTP CONFIG ----------------
	smtpCfg := config.LoadSMTPConfig()

	// ---------------- TWILIO CONFIG ----------------
	twilioCfg := config.LoadTwilioConfig()

	// ---------------- PROVIDERS ----------------
	smtpProvider := smtpprovider.NewSMTPProvider(
		smtpCfg,
	)

	twilioProvider := twilioprovider.NewTwilioProvider(
		twilioCfg,
	)

	// ---------------- MESSAGE SERVICE ----------------
	messageService := services.NewMessageService(
		smtpProvider,
		twilioProvider,
	)

	// ---------------- HANDLERS ----------------
	emailHandler := handlers.NewEmailHandler(
		messageService,
	)

	smsHandler := handlers.NewSMSHandler(
		messageService,
	)

	// ---------------- ROUTES ----------------
	app.Post("/email/send", emailHandler.SendEmail)

	app.Post("/sms/send", smsHandler.SendSMS)

	// ---------------- PORT ----------------
	port := os.Getenv("PORT")

	if port == "" {
		port = "3002"
	}

	log.Println("Message service running on", port)

	log.Fatal(app.Listen(":" + port))
}
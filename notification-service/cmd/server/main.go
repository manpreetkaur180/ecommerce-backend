package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"notification-service/internal/handlers"
	"notification-service/internal/services"
)

func main() {

	app := fiber.New()

	// -------- SERVICES --------
	notificationService := services.NewNotificationService()

	// -------- HANDLERS --------
	notificationHandler := handlers.NewNotificationHandler(
		notificationService,
	)

	// -------- ROUTES --------
	app.Post(
		"/notify/user-registered",
		notificationHandler.UserRegistered,
	)

	// -------- HEALTH CHECK --------
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// -------- PORT --------
	port := os.Getenv("PORT")

	if port == "" {
		port = "3003"
	}

	log.Println("Notification service running on", port)

	log.Fatal(app.Listen(":" + port))
}
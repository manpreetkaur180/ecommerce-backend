package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"user-service/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	godotenv.Load()
	app := fiber.New()

	// connect DB
	db := config.ConnectDB()
	_ = db // for now

	// health route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Println("Server running on port", port)
	log.Fatal(app.Listen(":" + port))
}
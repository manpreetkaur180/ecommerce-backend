package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"user-service/config"
	"user-service/internal/user"
)

func main() {
	godotenv.Load()

	app := fiber.New()

	// connect DB
	db := config.ConnectDB()

	// migrate table
	db.AutoMigrate(&user.User{}, &user.EmailVerification{})

	// redis
	rdb := config.ConnectRedis()

	// init layers
	userService := user.NewService(db, rdb)
	userHandler := user.NewHandler(userService)

	// register routes
	user.RegisterRoutes(app, userHandler)

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
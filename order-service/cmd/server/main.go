package main

import (
	"log"
	"os"

	"order-service/config"
	"order-service/internal/order"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// // JWT CHECK
	// if os.Getenv("JWT_SECRET") == "" {
	// 	log.Fatal("JWT_SECRET is required")
	// }

	// FIBER
	app := fiber.New()

	// DB CONNECTION
	db, err := config.ConnectDB()

	if err != nil {
		log.Fatal("failed to connect database")
	}

	// AUTO MIGRATION
	db.AutoMigrate(
	&order.Order{},
	&order.OrderItem{},
)

repo := order.NewRepository(db)

service := order.NewService(repo)

handler := order.NewHandler(service)

order.RegisterRoutes(app, handler)

	// PORT
	port := os.Getenv("PORT")

	if port == "" {
		port = "3006"
	}

	log.Println("Order service running on port", port)

	log.Fatal(app.Listen(":" + port))
}
package main

import (
	"log"
	"os"

	"cart-service/config"
	"cart-service/internal/cart"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is required")
	}

	app := fiber.New()

	// DB
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("DB connection failed")
	}

	// auto migrate
	db.AutoMigrate(&cart.Cart{}, &cart.CartItem{})

	// product client
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://product-service:3004"
	}

	productClient := cart.NewProductClient(productServiceURL)

	// service
	service := &cart.Service{
		DB:            db,
		ProductClient: productClient,
	}

	// handler
	handler := &cart.Handler{
		Service: service,
	}

	// routes
	cart.RegisterRoutes(app, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
	}

	log.Println("Cart service running on port", port)
	log.Fatal(app.Listen(":" + port))
}

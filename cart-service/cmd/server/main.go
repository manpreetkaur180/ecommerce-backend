package main

import (
	"log"
	"os"
	"cart-service/config"
	"cart-service/internal/cart"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

	// MIGRATE
	migrateDB(db)

	// product client
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://product-service:3004"
	}

	productClient := cart.NewProductClient(productServiceURL)

	// REPOSITORY
	repo := cart.NewRepository(db)

	// SERVICE
	service := cart.NewService(repo, productClient)

	// HANDLER
	handler := cart.NewHandler(service)

	// ROUTES
	cart.RegisterRoutes(app, handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
	}

	log.Println("Cart service running on port", port)
	log.Fatal(app.Listen(":" + port))
}

func migrateDB(db *gorm.DB) {
	if err := db.AutoMigrate(
		&cart.Cart{},
		&cart.CartItem{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migrated successfully")
}

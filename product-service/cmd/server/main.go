package main

import (
	"log"
	"os"

	"product-service/config"
	"product-service/internal/product"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is required")
	}

	app := fiber.New()

	// database
	db := config.ConnectDB()

	// migrate
	db.AutoMigrate(
		&product.Product{},
	)

	// initialize layers
	repo := product.NewRepository(db)
	productService := product.NewService(repo)
	productHandler := product.NewHandler(productService)

	// routes
	product.RegisterRoutes(app, productHandler)

	// health
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3004"
	}

	log.Println("Product service running on port", port)

	log.Fatal(app.Listen(":" + port))
}

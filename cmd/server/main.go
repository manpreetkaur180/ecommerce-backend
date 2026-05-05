package main

import (
	"ecommerce-backend/internal/database"
	"ecommerce-backend/internal/product"
	"ecommerce-backend/internal/user"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	db := database.Connect()

	db.AutoMigrate(
		&user.User{},
		&product.Product{},
	)

	//USERS
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	user.RegisterRoutes(app, userHandler)

	//PRODUC    TS

	productRepo := product.NewRepository(db)
	productService := product.NewService(productRepo)
	productHandler := product.NewHandler(productService)

	product.RegisterRoutes(app, productHandler)

	app.Listen(":3000")
} // package main

// import "github.com/gofiber/fiber/v2"

// func main() {
//     app := fiber.New()

//     app.Get("/", func(c *fiber.Ctx) error {
//         return c.SendString("API is running")
//     })

//     app.Listen(":3000")
// }

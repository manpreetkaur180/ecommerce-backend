package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"user-service/internal/seller"
		
	"user-service/config"
	"user-service/internal/user"
	"user-service/pkg/middleware"
)

func main() {
	godotenv.Load()

	app := fiber.New()
	app.Get(
	"/protected",
	middleware.RequireAuth(),
	func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"user_id":   c.Locals("user_id"),
			"role":      c.Locals("role"),
			"is_seller": c.Locals("is_seller"),
		})
	},
)
app.Get(
	"/seller-test",
	middleware.RequireAuth(),
	middleware.RequireSeller(),
	func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"message": "seller access granted",
		})
	},
)
app.Get(
	"/admin-test",
	middleware.RequireAuth(),
	middleware.RequireAdmin(),
	func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"message": "admin access granted",
		})
	},
)

	// connect DB
	db := config.ConnectDB()

db.AutoMigrate(
	&user.User{},
	&user.EmailVerification{},
	&user.PasswordResetVerification{},
	&user.PasswordReset{},
	&user.PasswordUpdate{},
	&seller.SellerApplication{},
)
user.SeedAdmin(db)

	// redis
	rdb := config.ConnectRedis()

	// init layers
	userService := user.NewService(db, rdb)
	userHandler := user.NewHandler(userService)
	sellerService := seller.NewService(db)
sellerHandler := seller.NewHandler(sellerService)

	// register routes
	user.RegisterRoutes(app, userHandler)
	seller.RegisterRoutes(app, sellerHandler)

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

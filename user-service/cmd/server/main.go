package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"user-service/internal/seller"

	"gorm.io/gorm"
	"user-service/config"
	"user-service/internal/user"
	"user-service/pkg/middleware"
	"user-service/pkg/utils"
)

func main() {
	godotenv.Load()

	if err := utils.ValidateJWTSecret(); err != nil {
		log.Fatal(err)
	}

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
		middleware.RequireRoles("admin"),
		func(c *fiber.Ctx) error {

			return c.JSON(fiber.Map{
				"message": "admin access granted",
			})
		},
	)

	// connect DB
	db := config.ConnectDB()

	dropLegacySellerApplicationUniqueIndex(db)

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

func dropLegacySellerApplicationUniqueIndex(db *gorm.DB) {
	const indexName = "idx_seller_applications_user_id"

	if !db.Migrator().HasIndex(&seller.SellerApplication{}, indexName) {
		return
	}

	var isUnique bool
	if err := db.Raw(`
		SELECT i.indisunique
		FROM pg_class c
		JOIN pg_index i ON i.indexrelid = c.oid
		WHERE c.relname = ?
		LIMIT 1
	`, indexName).Scan(&isUnique).Error; err != nil {
		log.Println("Failed to inspect seller application user index:", err)
		return
	}

	if !isUnique {
		return
	}

	if err := db.Migrator().DropIndex(&seller.SellerApplication{}, indexName); err != nil {
		log.Println("Failed to drop old seller application user unique index:", err)
	}
}

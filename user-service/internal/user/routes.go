package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RegisterRoutes(app *fiber.App, handler *Handler) {
	api := app.Group("/api/v1")
	userRoutes := api.Group("/user")

	// -------- RATE LIMIT CONFIG --------
	authLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: time.Minute,

		// Custom response
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many requests",
				"code":  429,
			})
		},
	})

	// Apply limiter ONLY to sensitive routes
	userRoutes.Post("/register", authLimiter, handler.Register)
	userRoutes.Post("/login", authLimiter, handler.Login)
}
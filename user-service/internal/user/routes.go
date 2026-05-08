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
		Max:        5,
		Expiration: time.Minute,

		// Custom response
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many requests",
				"code":  429,
			})
		},
	})


	userRoutes.Post("/register", authLimiter, handler.Register)
	userRoutes.Post("/login", authLimiter, handler.Login)
	userRoutes.Post("/send-otp", authLimiter, handler.SendOTP)
	userRoutes.Post("/login-otp", authLimiter, handler.LoginWithOTP)
	userRoutes.Get("/verify-email",handler.VerifyEmail)

}
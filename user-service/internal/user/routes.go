package user

import (
	"time"
	"user-service/pkg/middleware"

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

		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many requests",
				"code":  429,
			})
		},
	})
	userRoutes.Post("/address",middleware.RequireAuth(),handler.AddAddress)
	userRoutes.Get("/address",middleware.RequireAuth(),handler.GetAddresses	)
	userRoutes.Post("/register", authLimiter, handler.Register)
	userRoutes.Post("/login", authLimiter, handler.Login)
	userRoutes.Post("/refresh-token", middleware.RequireAuth(), handler.RefreshToken)
	userRoutes.Post("/send-otp", authLimiter, handler.SendOTP)
	userRoutes.Post("/login-otp", authLimiter, handler.LoginWithOTP)
	userRoutes.Get("/verify-email", handler.VerifyEmail)
	userRoutes.Post("/forgot-password", authLimiter, handler.ForgotPassword)
	userRoutes.Get("/reset-password", handler.ResetPasswordForm)
	userRoutes.Post("/reset-password", handler.ResetPassword)
	userRoutes.Post("/update-password", authLimiter, handler.RequestUpdatePassword)
	userRoutes.Get("/update-password", handler.UpdatePasswordForm)
	userRoutes.Post("/update-password/confirm", handler.ConfirmUpdatePassword)
}

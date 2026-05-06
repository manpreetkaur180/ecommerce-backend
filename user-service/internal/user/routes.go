package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *Handler) {
	api := app.Group("/api/v1")
	userRoutes := api.Group("/user")

	userRoutes.Post("/register", handler.Register)
	userRoutes.Post("/login", handler.Login)
}
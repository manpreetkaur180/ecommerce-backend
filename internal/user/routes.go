package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *Handler) {
    user := app.Group("/users")

    user.Post("/", handler.CreateUser)
    user.Post("/login", handler.Login)
}
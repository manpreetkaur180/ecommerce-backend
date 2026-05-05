package order

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *Handler) {
    order := app.Group("/orders")

    order.Post("/", handler.CreateOrder)
}
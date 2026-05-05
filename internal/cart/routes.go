package cart

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *Handler) {
    cart := app.Group("/cart")

    cart.Post("/", handler.AddToCart)
}
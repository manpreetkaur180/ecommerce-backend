package product

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *Handler) {
    product := app.Group("/products")

    product.Post("/", handler.CreateProduct)
    product.Get("/", handler.GetProducts)
}
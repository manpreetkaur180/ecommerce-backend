package order

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler) {

	api := app.Group("/api/v1")

	orders := api.Group("/orders")

	orders.Post("/", h.CreateOrder)
	orders.Get("/", h.GetMyOrders)

	seller := api.Group("/seller/orders")

	seller.Patch("/:id/confirm", h.ConfirmOrder)
	seller.Patch("/:id/reject", h.RejectOrder)
}
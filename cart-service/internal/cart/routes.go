package cart

import (
	"cart-service/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *Handler) {

	api := app.Group("/api/v1")

	// -----------------------------
	// CART ROUTES (BUYER ONLY)
	// -----------------------------
	cart := api.Group(
		"/cart",
		middleware.RequireAuth(),
		middleware.RequireRoles("buyer"),
	)

	// GET CART
	cart.Get("/", handler.GetCart)

	// ADD TO CART
	cart.Post("/add", handler.AddToCart)

	// REDUCE QUANTITY
	cart.Patch("/reduce", handler.ReduceItem)

	// REMOVE ITEM
	cart.Delete("/item/:product_id", handler.RemoveItem)

	// CLEAR CART
	cart.Delete("/", handler.ClearCart)
}
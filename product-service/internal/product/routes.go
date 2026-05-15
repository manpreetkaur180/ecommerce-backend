package product

import (
	"product-service/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	handler *Handler,
) {

	api := app.Group("/api/v1")

	// -----------------------------
	// PUBLIC BUYER ROUTES
	// -----------------------------
	buyer := api.Group(
		"/buyer",
		middleware.RequireAuth(),
		middleware.RequireRoles("buyer"),
	)

	buyer.Get(
		"/products",
		handler.GetAllProducts,
	)

	buyer.Get(
		"/products/:id",
		handler.GetProductByID,
	)
	buyer.Post("/products/bulk", handler.GetProductsByIDs)

	// -----------------------------
	// SELLER ROUTES
	// -----------------------------

	seller := api.Group(
		"/seller",
		middleware.RequireAuth(),
		middleware.RequireRoles("seller"),
	)

	seller.Post(
		"/products",
		handler.CreateProduct,
	)

	seller.Get(
		"/products",
		handler.GetSellerProducts,
	)

	seller.Get(
		"/products/:id",
		handler.GetProductByID,
	)

	seller.Patch(
		"/products/:id",
		handler.UpdateProduct,
	)

	seller.Delete(
		"/products/:id",
		handler.DeleteProduct,
	)
	internal := api.Group("/internal/products")

internal.Get("/:id/inventory", handler.GetInventory)

internal.Patch("/:id/reduce-stock", handler.ReduceStock)
}

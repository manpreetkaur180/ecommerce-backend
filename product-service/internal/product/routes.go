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
		middleware.RequireBuyer(),
	)

	buyer.Get(
		"/products",
		handler.GetAllProducts,
	)

	buyer.Get(
		"/products/:id",
		handler.GetProductByID,
	)

	// -----------------------------
	// SELLER ROUTES
	// -----------------------------

	seller := api.Group(
		"/seller",
		middleware.RequireAuth(),
		middleware.RequireSeller(),
	)

	seller.Post(
		"/products",
		handler.CreateProduct,
	)

	seller.Get(
		"/products",
		handler.GetSellerProducts,
	)

	seller.Patch(
		"/products/:id",
		handler.UpdateProduct,
	)

	seller.Delete(
		"/products/:id",
		handler.DeleteProduct,
	)
}

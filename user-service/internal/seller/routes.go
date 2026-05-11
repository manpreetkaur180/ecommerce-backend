package seller

import (
	"user-service/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	handler *Handler,
) {

	seller := app.Group(
		"/api/v1/seller",
		middleware.RequireAuth(),
	)

	seller.Post(
		"/apply",
		handler.ApplySeller,
	)
}
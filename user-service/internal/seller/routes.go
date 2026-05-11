package seller

import (
	"user-service/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	handler *Handler,
) {

	// buyer routes
	seller := app.Group(
		"/api/v1/seller",
		middleware.RequireAuth(),
	)

	seller.Post(
		"/apply",
		handler.ApplySeller,
	)

	// admin routes
	admin := app.Group(
		"/api/v1/admin",
		middleware.RequireAuth(),
		middleware.RequireAdmin(),
	)

	admin.Get(
		"/seller-applications",
		handler.GetAllApplications,
	)

	admin.Patch(
		"/seller-applications/:id/approve",
		handler.ApproveApplication,
	)

	admin.Patch(
		"/seller-applications/:id/reject",
		handler.RejectApplication,
	)
}
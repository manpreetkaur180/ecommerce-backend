package seller

import (
	"user-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) ApplySeller(c *fiber.Ctx) error {

	var req ApplySellerRequest

	// parse body
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid request body",
		)
	}

	// get authenticated user id from middleware
	userID, ok := c.Locals("user_id").(uint)

	if !ok {
		return utils.ErrorResponse(
			c,
			401,
			"unauthorized",
		)
	}

	// call service
	if err := h.Service.ApplySeller(userID, req); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"seller application submitted successfully",
		nil,
	)
}
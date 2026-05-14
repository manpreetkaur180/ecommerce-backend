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
func (h *Handler) GetAllApplications(c *fiber.Ctx) error {

	applications, err := h.Service.GetAllApplications()

	if err != nil {
		return utils.ErrorResponse(
			c,
			500,
			err.Error(),
		)
	}

	response := make([]SellerApplicationListResponse, 0, len(applications))
	for _, application := range applications {
		response = append(response, SellerApplicationListResponse{
			ID:                  application.ID,
			UserID:              application.UserID,
			BusinessName:        application.BusinessName,
			BusinessDescription: application.BusinessDescription,
			GSTIN:               application.GSTIN,
			AadharNumber:        application.AadharNumber,
			Status:              application.Status,
			CreatedAt:           application.CreatedAt,
			UpdatedAt:           application.UpdatedAt,
		})
	}

	return utils.SuccessResponse(
		c,
		200,
		"seller applications fetched successfully",
		response,
	)
}
func (h *Handler) ApproveApplication(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid application id",
		)
	}

	var req SellerApplicationDecisionRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid request body",
		)
	}

	application, err := h.Service.ApproveApplication(
		uint(id),
		req.AdminNote,
	)
	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"seller application approved successfully",
		ToSellerApplicationDecisionResponse(application),
	)
}
func (h *Handler) RejectApplication(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid application id",
		)
	}

	var req SellerApplicationDecisionRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid request body",
		)
	}

	application, err := h.Service.RejectApplication(
		uint(id),
		req.AdminNote,
	)
	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"seller application rejected successfully",
		ToSellerApplicationDecisionResponse(application),
	)
}

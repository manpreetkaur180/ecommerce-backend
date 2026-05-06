package user

import (
	"user-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

// -------- REGISTER --------
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	// parse request
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	// call service
	user, err := h.Service.Register(req)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// success response
	return utils.SuccessResponse(
		c,
		200,
		"Hi "+user.Name+", successfully registered",
		nil,
	)
}

// -------- LOGIN --------
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	// parse request
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	// call service
	user, err := h.Service.Login(req)
	if err != nil {
		return utils.ErrorResponse(c, 401, err.Error())
	}

	// success response
	return utils.SuccessResponse(
		c,
		200,
		"Hi "+user.Name+", logged in successfully",
		nil,
	)
}
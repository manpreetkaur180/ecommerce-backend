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
func (h *Handler) SendOTP(c *fiber.Ctx) error {

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	req.Email = utils.NormalizeEmail(req.Email)
	req.Phone = utils.NormalizePhone(req.Phone)

	if req.Email == "" && req.Phone == "" {
		return utils.ErrorResponse(c, 400, "email or phone is required")
	}

	if req.Email != "" {
		if err := utils.ValidateEmail(req.Email); err != nil {
			return utils.ErrorResponse(c, 400, err.Error())
		}
	}

	if req.Phone != "" {
		if err := utils.ValidatePhone(req.Phone); err != nil {
			return utils.ErrorResponse(c, 400, err.Error())
		}
	}

	user, err := h.Service.FindByIdentifier(
		req.Email,
		req.Phone,
	)

	if err != nil {
		return utils.ErrorResponse(c, 400, "user not found")
	}

	channel := "email"

	identifier := req.Email

	if req.Phone != "" {
		channel = "phone"
		identifier = req.Phone
	}

	err = h.Service.SendOTP(
		user,
		identifier,
		channel,
	)

	if err != nil {
		return utils.ErrorResponse(c, 500, "failed to send otp")
	}

	return utils.SuccessResponse(
		c,
		200,
		"OTP sent successfully",
		nil,
	)
}
func (h *Handler) LoginWithOTP(c *fiber.Ctx) error {
	var req OTPLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	req.Email = utils.NormalizeEmail(req.Email)
	req.Phone = utils.NormalizePhone(req.Phone)

	if req.Email == "" && req.Phone == "" {
		return utils.ErrorResponse(c, 400, "email or phone is required")
	}

	if req.OTP == "" {
		return utils.ErrorResponse(c, 400, "otp is required")
	}

	identifier := req.Email

	if identifier == "" {
		identifier = req.Phone
	}

	if err := h.Service.VerifyOTP(identifier, req.OTP); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	user, err := h.Service.FindByIdentifier(req.Email, req.Phone)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid request")
	}

	return utils.SuccessResponse(
		c,
		200,
		"Hi "+user.Name+", logged in successfully",
		nil,
	)
}


func (h *Handler) VerifyEmail(c *fiber.Ctx) error {

	token := c.Query("token")

	if token == "" {
		return utils.ErrorResponse(c, 400, "token required")
	}

	err := h.Service.VerifyEmail(token)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(
		c,
		200,
		"email verified successfully",
		nil,
	)
}
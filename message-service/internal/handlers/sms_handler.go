package handlers

import (

	"message-service/internal/models"
	"message-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SMSHandler struct {
	service *services.MessageService
}

func NewSMSHandler(
	service *services.MessageService,
) *SMSHandler {

	return &SMSHandler{
		service: service,
	}
}

func (h *SMSHandler) SendSMS(
	c *fiber.Ctx,
) error {

	var req models.SMSRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	err := h.service.SendSMS(req)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "sms sent successfully",
	})
}
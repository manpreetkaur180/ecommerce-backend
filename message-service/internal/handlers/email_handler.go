package handlers

import (
	"message-service/internal/models"
	"message-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type EmailHandler struct {
	service *services.MessageService
}

func NewEmailHandler(service *services.MessageService) *EmailHandler {
	return &EmailHandler{
		service: service,
	}
}

func (h *EmailHandler) SendEmail(c *fiber.Ctx) error {

	var req models.EmailRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	err := h.service.SendEmail(req)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "email sent successfully",
	})
}
package handlers


import (
	"notification-service/internal/models"
	"notification-service/internal/services"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(
	service *services.NotificationService,
) *NotificationHandler {

	return &NotificationHandler{
		service: service,
	}
}

func (h *NotificationHandler) UserRegistered(
	c *fiber.Ctx,
) error {

	var event models.UserRegisteredEvent

	if err := c.BodyParser(&event); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	err := h.service.HandleUserRegistered(
		event.Name,
		event.Email,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "notification processed",
	})
}

package user

import "github.com/gofiber/fiber/v2"

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service}
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
    var user User

    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(err)
    }

    if err := h.service.CreateUser(&user); err != nil {
        return c.Status(500).JSON(err)
    }

    return c.JSON(user)
}
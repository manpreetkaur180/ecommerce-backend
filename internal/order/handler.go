package order

import "github.com/gofiber/fiber/v2"

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service}
}

type CreateOrderRequest struct {
    UserID uint `json:"user_id"`
}

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
    var req CreateOrderRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(err)
    }

    err := h.service.CreateOrder(req.UserID)
    if err != nil {
        return c.Status(500).JSON(err)
    }

    return c.JSON("order created")
}
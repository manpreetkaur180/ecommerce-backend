package cart

import (


    "github.com/gofiber/fiber/v2"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service}
}

type AddRequest struct {
    UserID    uint `json:"user_id"`
    ProductID uint `json:"product_id"`
    Quantity  int  `json:"quantity"`
}

func (h *Handler) AddToCart(c *fiber.Ctx) error {
    var req AddRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(err)
    }

    err := h.service.AddToCart(req.UserID, req.ProductID, req.Quantity)
    if err != nil {
        return c.Status(500).JSON(err)
    }

    return c.JSON("added to cart")
}
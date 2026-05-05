package product

import "github.com/gofiber/fiber/v2"

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service}
}

func (h *Handler) CreateProduct(c *fiber.Ctx) error {
    var product Product

    if err := c.BodyParser(&product); err != nil {
        return c.Status(400).JSON(err)
    }

    if err := h.service.CreateProduct(&product); err != nil {
        return c.Status(500).JSON(err)
    }

    return c.JSON(product)
}

func (h *Handler) GetProducts(c *fiber.Ctx) error {
    products, err := h.service.GetProducts()
    if err != nil {
        return c.Status(500).JSON(err)
    }

    return c.JSON(products)
}
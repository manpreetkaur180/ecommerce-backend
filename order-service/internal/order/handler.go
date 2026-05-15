package order

import (
	"order-service/pkg/utils"

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

func (h *Handler) CreateOrder(c *fiber.Ctx) error {

	userID := c.Locals("user_id").(uint)

	token := c.Get("Authorization")

	var req CreateOrderRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid body")
	}

	order, err := h.Service.CreateOrder(
		userID,
		token,
		req,
	)

	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(
		c,
		200,
		"order created successfully",
		order,
	)
}

func (h *Handler) GetMyOrders(c *fiber.Ctx) error {

	userID := c.Locals("user_id").(uint)

	orders, err := h.Service.GetUserOrders(userID)
	if err != nil {
		return utils.ErrorResponse(c, 500, "failed to fetch orders")
	}

	return utils.SuccessResponse(
		c,
		200,
		"orders fetched successfully",
		orders,
	)
}

func (h *Handler) ConfirmOrder(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) RejectOrder(c *fiber.Ctx) error {
	return nil
}
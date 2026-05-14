package cart

import (
	"cart-service/pkg/utils"

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

func (h *Handler) AddToCart(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	var req AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}
	if req.ProductID == 0 {
		return utils.ErrorResponse(c, 400, "product id is required")
	}

	if req.Quantity < 1 {
		return utils.ErrorResponse(c, 400, "quantity must be at least 1")
	}

	cart, err := h.Service.AddToCart(
		userID,
		req.ProductID,
		req.Quantity,
		c.Get("Authorization"),
	)

	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(
		c,
		200,
		"product added to cart successfully",
		cart,
	)
}

func (h *Handler) GetCart(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	cart, err := h.Service.GetCart(userID, c.Get("Authorization"))
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 200, "cart fetched successfully", cart)
}

func (h *Handler) ReduceItem(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	var req ReduceCartRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}
	if req.ProductID == 0 {
		return utils.ErrorResponse(c, 400, "product id is required")
	}
	err := h.Service.ReduceItem(
		userID,
		req.ProductID,
	)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 200, "cart item quantity reduced successfully", nil)
}

func (h *Handler) RemoveItem(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	productID, err := c.ParamsInt("product_id")
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid product id")
	}

	err = h.Service.RemoveItem(userID, uint(productID))
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 200, "cart item removed successfully", nil)
}

func (h *Handler) ClearCart(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	err := h.Service.ClearCart(userID)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 200, "cart cleared successfully", nil)
}

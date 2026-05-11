package product

import (
	"product-service/pkg/utils"

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
func (h *Handler) CreateProduct(c *fiber.Ctx) error {

	var req CreateProductRequest

	// parse request
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid request body",
		)
	}

	// get seller id from jwt middleware
	sellerID, ok := c.Locals("user_id").(uint)

	if !ok {
		return utils.ErrorResponse(
			c,
			401,
			"unauthorized",
		)
	}

	product, err := h.Service.CreateProduct(
		sellerID,
		req,
	)

	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"product created successfully",
		product,
	)
}
func (h *Handler) GetAllProducts(c *fiber.Ctx) error {

	products, err := h.Service.GetAllProducts()

	if err != nil {
		return utils.ErrorResponse(
			c,
			500,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"products fetched successfully",
		products,
	)
}
func (h *Handler) GetProductByID(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid product id",
		)
	}

	product, err := h.Service.GetProductByID(uint(id))

	if err != nil {
		return utils.ErrorResponse(
			c,
			404,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"product fetched successfully",
		product,
	)
}
func (h *Handler) GetSellerProducts(c *fiber.Ctx) error {

	sellerID, ok := c.Locals("user_id").(uint)

	if !ok {
		return utils.ErrorResponse(
			c,
			401,
			"unauthorized",
		)
	}

	products, err := h.Service.GetSellerProducts(sellerID)

	if err != nil {
		return utils.ErrorResponse(
			c,
			500,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"seller products fetched successfully",
		products,
	)
}
func (h *Handler) UpdateProduct(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid product id",
		)
	}

	sellerID, ok := c.Locals("user_id").(uint)

	if !ok {
		return utils.ErrorResponse(
			c,
			401,
			"unauthorized",
		)
	}

	var req UpdateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid request body",
		)
	}

	product, err := h.Service.UpdateProduct(
		sellerID,
		uint(id),
		req,
	)

	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"product updated successfully",
		product,
	)
}
func (h *Handler) DeleteProduct(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")

	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid product id",
		)
	}

	sellerID, ok := c.Locals("user_id").(uint)

	if !ok {
		return utils.ErrorResponse(
			c,
			401,
			"unauthorized",
		)
	}

	if err := h.Service.DeleteProduct(
		sellerID,
		uint(id),
	); err != nil {

		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"product deleted successfully",
		nil,
	)
}
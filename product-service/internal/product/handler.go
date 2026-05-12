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

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid request body",
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

	page := c.QueryInt("page", 1)

	if page < 1 {
		page = 1
	}

	products, err := h.Service.GetAllProducts(page)

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

	if err := utils.ValidateID(id, "product id"); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
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

	page := c.QueryInt("page", 1)

	if page < 1 {
		page = 1
	}

	products, err := h.Service.GetSellerProducts(
		sellerID,
		page,
	)

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

	if err := utils.ValidateID(id, "product id"); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
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

	if err := utils.ValidateID(id, "product id"); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
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
		"product deactivated successfully",
		nil,
	)
}
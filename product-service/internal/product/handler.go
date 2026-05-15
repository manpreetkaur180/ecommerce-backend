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

func validateAndNormalizeCreateProduct(req *CreateProductRequest) error {
	req.Title = utils.NormalizeTitle(req.Title)
	req.Description = utils.NormalizeDescription(req.Description)
	req.Category = utils.NormalizeCategory(req.Category)
	req.ImageURLs = utils.NormalizeStringSlice(req.ImageURLs)
	req.Offers = utils.NormalizeOptionalText(req.Offers)
	req.Warranty = utils.NormalizeOptionalText(req.Warranty)

	if err := utils.ValidateTitle(req.Title); err != nil {
		return err
	}
	if err := utils.ValidateDescription(req.Description); err != nil {
		return err
	}
	if err := utils.ValidatePrice(req.Price); err != nil {
		return err
	}
	if err := utils.ValidateStock(req.Stock); err != nil {
		return err
	}
	if err := utils.ValidateCategory(req.Category); err != nil {
		return err
	}
	if err := utils.ValidateImageURLs(req.ImageURLs); err != nil {
		return err
	}
	if err := utils.ValidateOffers(req.Offers); err != nil {
		return err
	}
	if err := utils.ValidateWarranty(req.Warranty); err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetProductsByIDs(c *fiber.Ctx) error {
	var req BulkProductIDsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	if len(req.ProductIDs) == 0 {
		return utils.ErrorResponse(c, 400, "product ids are required")
	}

	products, err := h.Service.GetProductsByIDs(req.ProductIDs)
	if err != nil {
		return utils.ErrorResponse(c, 500, err.Error())
	}

	return utils.SuccessResponse(c, 200, "products fetched successfully", products)
}
func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	if err := validateAndNormalizeCreateProduct(&req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	sellerID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	product, err := h.Service.CreateProduct(sellerID, req)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 200, "product created successfully", product)
}

func (h *Handler) GetAllProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	products, err := h.Service.GetAllProducts(page)
	if err != nil {
		return utils.ErrorResponse(c, 500, err.Error())
	}

	return utils.SuccessResponse(c, 200, "products fetched successfully", products)
}

func (h *Handler) GetProductByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid product id")
	}

	if err := utils.ValidateID(id, "product id"); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	role, ok := c.Locals("role").(string)
	if !ok || role == "" {
		return utils.ErrorResponse(c, 403, "role not found")
	}

	if role == "seller" {
		sellerID, ok := c.Locals("user_id").(uint)
		if !ok {
			return utils.ErrorResponse(c, 401, "unauthorized")
		}

		product, err := h.Service.GetSellerProductByID(sellerID, uint(id))
		if err != nil {
			return utils.ErrorResponse(c, 404, err.Error())
		}

		return utils.SuccessResponse(c, 200, "seller product fetched successfully", product)
	}

	if role != "buyer" {
		return utils.ErrorResponse(c, 403, "access denied")
	}

	product, err := h.Service.GetProductByID(uint(id))
	if err != nil {
		return utils.ErrorResponse(c, 404, err.Error())
	}

	return utils.SuccessResponse(c, 200, "product fetched successfully", product)
}

func (h *Handler) GetSellerProducts(c *fiber.Ctx) error {
	sellerID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	products, err := h.Service.GetSellerProducts(sellerID, page)
	if err != nil {
		return utils.ErrorResponse(c, 500, err.Error())
	}

	return utils.SuccessResponse(c, 200, "seller products fetched successfully", products)
}

func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid product id")
	}

	if err := utils.ValidateID(id, "product id"); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	sellerID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	product, err := h.Service.UpdateProduct(sellerID, uint(id), req)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 200, "product updated successfully", product)
}

func (h *Handler) DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid product id")
	}

	if err := utils.ValidateID(id, "product id"); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	sellerID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	if err := h.Service.DeleteProduct(sellerID, uint(id)); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 200, "product deactivated successfully", nil)
}
func (h *Handler) GetInventory(c *fiber.Ctx) error {

	productID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid product id")
	}

	inventory, err := h.Service.GetInventory(uint(productID))
	if err != nil {
		return utils.ErrorResponse(c, 404, err.Error())
	}

	return utils.SuccessResponse(
		c,
		200,
		"inventory fetched successfully",
		inventory,
	)
}
func (h *Handler) ReduceStock(c *fiber.Ctx) error {

	productID, err := c.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid product id")
	}

	type reqBody struct {
		Quantity int `json:"quantity"`
	}

	var req reqBody

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid body")
	}

	if err := h.Service.ReduceStock(
		uint(productID),
		req.Quantity,
	); err != nil {

		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(
		c,
		200,
		"stock updated successfully",
		nil,
	)
}
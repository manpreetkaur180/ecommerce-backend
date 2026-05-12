package product

import (
	"errors"
	"product-service/pkg/utils"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

func (s *Service) CreateProduct(
	sellerID uint,
	req CreateProductRequest,
) (*Product, error) {

	// -------------------------
	// NORMALIZATION
	// -------------------------

	req.Title = utils.NormalizeTitle(req.Title)
	req.Description = utils.NormalizeDescription(req.Description)
	req.Category = utils.NormalizeCategory(req.Category)
	req.ImageURLs = utils.NormalizeStringSlice(req.ImageURLs)
	req.Offers = utils.NormalizeOptionalText(req.Offers)
	req.Warranty = utils.NormalizeOptionalText(req.Warranty)

	// -------------------------
	// VALIDATIONS
	// -------------------------

	if err := utils.ValidateTitle(req.Title); err != nil {
		return nil, err
	}

	if err := utils.ValidateDescription(req.Description); err != nil {
		return nil, err
	}

	if err := utils.ValidatePrice(req.Price); err != nil {
		return nil, err
	}

	if err := utils.ValidateStock(req.Stock); err != nil {
		return nil, err
	}

	if err := utils.ValidateCategory(req.Category); err != nil {
		return nil, err
	}

	if err := utils.ValidateImageURLs(req.ImageURLs); err != nil {
		return nil, err
	}

	if err := utils.ValidateOffers(req.Offers); err != nil {
		return nil, err
	}

	if err := utils.ValidateWarranty(req.Warranty); err != nil {
		return nil, err
	}

	// -------------------------
	// CREATE PRODUCT
	// -------------------------

	product := Product{
		SellerID: sellerID,

		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,

		ImageURLs: datatypes.JSONSlice[string](req.ImageURLs),

		Offers: req.Offers,

		ReturnAvailable: req.ReturnAvailable,

		Warranty: req.Warranty,

		IsActive: true,
	}

	if err := s.DB.Create(&product).Error; err != nil {
		return nil, errors.New("failed to create product")
	}

	return &product, nil
}

func (s *Service) GetAllProducts() ([]BuyerProductResponse, error) {

	var products []Product

	if err := s.DB.Where(
		"is_active = ?",
		true,
	).Order("created_at desc").
		Find(&products).Error; err != nil {

		return nil, errors.New("failed to fetch products")
	}

	var response []BuyerProductResponse

	for _, product := range products {

		image := ""

		if len(product.ImageURLs) > 0 {
			image = product.ImageURLs[0]
		}

		response = append(response, BuyerProductResponse{
			ID:          product.ID,
			Title:       product.Title,
			ImageURL:    image,
			Description: product.Description,
			Price:       product.Price,
			Offers:      product.Offers,

			ExpectedDelivery: utils.GetExpectedDelivery(),
		})
	}

	return response, nil
}

func (s *Service) GetProductByID(
	productID uint,
) (*BuyerProductDetailResponse, error) {

	var product Product

	if err := s.DB.Where(
		"id = ? AND is_active = ?",
		productID,
		true,
	).First(&product).Error; err != nil {

		return nil, errors.New("product not found")
	}

	response := BuyerProductDetailResponse{
		ID:              product.ID,
		Title:           product.Title,
		Description:     product.Description,
		Price:           product.Price,
		Category:        product.Category,
		ImageURLs:       product.ImageURLs,
		Offers:          product.Offers,
		Rating:          product.Rating,
		ReturnAvailable: product.ReturnAvailable,
		Warranty:        product.Warranty,
		InStock:         product.Stock > 0,

		ExpectedDelivery: utils.GetExpectedDelivery(),
	}

	return &response, nil
}

func (s *Service) GetSellerProducts(
	sellerID uint,
) ([]Product, error) {

	var products []Product

	if err := s.DB.Where(
		"seller_id = ?",
		sellerID,
	).Order("created_at desc").
		Find(&products).Error; err != nil {

		return nil, errors.New("failed to fetch seller products")
	}

	return products, nil
}

func (s *Service) UpdateProduct(
	sellerID uint,
	productID uint,
	req UpdateProductRequest,
) (*Product, error) {

	var product Product

	// ownership validation
	if err := s.DB.Where(
		"id = ? AND seller_id = ?",
		productID,
		sellerID,
	).First(&product).Error; err != nil {

		return nil, errors.New("product not found")
	}

	// -------------------------
	// UPDATE FIELDS
	// -------------------------

	if req.Title != nil {

		title := utils.NormalizeTitle(*req.Title)

		if err := utils.ValidateTitle(title); err != nil {
			return nil, err
		}

		product.Title = title
	}

	if req.Description != nil {

		description := utils.NormalizeDescription(*req.Description)

		if err := utils.ValidateDescription(description); err != nil {
			return nil, err
		}

		product.Description = description
	}

	if req.Price != nil {

		if err := utils.ValidatePrice(*req.Price); err != nil {
			return nil, err
		}

		product.Price = *req.Price
	}

	if req.Stock != nil {

		if err := utils.ValidateStock(*req.Stock); err != nil {
			return nil, err
		}

		product.Stock = *req.Stock
	}

	if req.Category != nil {

		category := utils.NormalizeCategory(*req.Category)

		if err := utils.ValidateCategory(category); err != nil {
			return nil, err
		}

		product.Category = category
	}

	if req.ImageURLs != nil {
		imageURLs := utils.NormalizeStringSlice(*req.ImageURLs)

		if err := utils.ValidateImageURLs(imageURLs); err != nil {
			return nil, err
		}

		product.ImageURLs = datatypes.JSONSlice[string](imageURLs)
	}

	if req.Offers != nil {

		offers := utils.NormalizeOptionalText(*req.Offers)

		if err := utils.ValidateOffers(offers); err != nil {
			return nil, err
		}

		product.Offers = offers
	}

	if req.ReturnAvailable != nil {
		product.ReturnAvailable = *req.ReturnAvailable
	}

	if req.Warranty != nil {

		warranty := utils.NormalizeOptionalText(*req.Warranty)

		if err := utils.ValidateWarranty(warranty); err != nil {
			return nil, err
		}

		product.Warranty = warranty
	}

	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.DB.Save(&product).Error; err != nil {
		return nil, errors.New("failed to update product")
	}

	return &product, nil
}

func (s *Service) DeleteProduct(
	sellerID uint,
	productID uint,
) error {

	var product Product

	// ownership validation
	if err := s.DB.Where(
		"id = ? AND seller_id = ?",
		productID,
		sellerID,
	).First(&product).Error; err != nil {

		return errors.New("product not found")
	}

	product.IsActive = false

	if err := s.DB.Save(&product).Error; err != nil {
		return errors.New("failed to deactivate product")
	}

	return nil
}

package product

import (
	"errors"

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

	// validations
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	if req.Price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}

	if req.Stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	// create product
	product := Product{
		SellerID:    sellerID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURLs:    req.ImageURLs,
		Offers: req.Offers,
		IsActive:    true,
		Rating: req.Rating,

ReturnAvailable: req.ReturnAvailable,

Warranty: req.Warranty,
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

			ExpectedDelivery: "16 May - 17 May",
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
		ID:                product.ID,
		Title:             product.Title,
		Description:       product.Description,
		Price:             product.Price,
		Category:          product.Category,
		ImageURLs:         product.ImageURLs,
		Offers:            product.Offers,
		Rating:            product.Rating,
		ReturnAvailable:   product.ReturnAvailable,
		Warranty:          product.Warranty,
		InStock:           product.Stock > 0,
		AvailableQuantity: product.Stock,

		ExpectedDelivery: "16 May - 17 May",
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

	// update fields
	if req.Title != "" {
		product.Title = req.Title
	}

	if req.Description != "" {
		product.Description = req.Description
	}

	if req.Price > 0 {
		product.Price = req.Price
	}

	if req.Stock != nil {
		if *req.Stock < 0 {
			return nil, errors.New("stock cannot be negative")
		}
		product.Stock = *req.Stock
	}

	if req.Category != "" {
		product.Category = req.Category
	}

	if len(req.ImageURLs) > 0 {
		product.ImageURLs = req.ImageURLs
	}
	if req.Offers != "" {
	product.Offers = req.Offers
}	
   if req.Rating != nil {
	product.Rating = *req.Rating
}

if req.ReturnAvailable != nil {
	product.ReturnAvailable = *req.ReturnAvailable
}

if req.Warranty != "" {
	product.Warranty = req.Warranty
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

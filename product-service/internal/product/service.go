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
		SellerID:   sellerID,
		Title:      req.Title,
		Description: req.Description,
		Price:      req.Price,
		Stock:      req.Stock,
		Category:   req.Category,
		ImageURL:   req.ImageURL,
		IsActive:   true,
	}

	if err := s.DB.Create(&product).Error; err != nil {
		return nil, errors.New("failed to create product")
	}

	return &product, nil
}
func (s *Service) GetAllProducts() ([]Product, error) {

	var products []Product

	if err := s.DB.Where(
		"is_active = ?",
		true,
	).Order("created_at desc").
		Find(&products).Error; err != nil {

		return nil, errors.New("failed to fetch products")
	}

	return products, nil
}	
func (s *Service) GetProductByID(
	productID uint,
) (*Product, error) {

	var product Product

	if err := s.DB.Where(
		"id = ? AND is_active = ?",
		productID,
		true,
	).First(&product).Error; err != nil {

		return nil, errors.New("product not found")
	}

	return &product, nil
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

	if req.Stock >= 0 {
		product.Stock = req.Stock
	}

	if req.Category != "" {
		product.Category = req.Category
	}

	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
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

	if err := s.DB.Delete(&product).Error; err != nil {
		return errors.New("failed to delete product")
	}

	return nil
}
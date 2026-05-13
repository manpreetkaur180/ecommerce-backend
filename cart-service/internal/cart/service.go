package cart

import (
	"errors"
	"cart-service/pkg/utils"
	"gorm.io/gorm"
)

type Service struct {
	DB            *gorm.DB
	ProductClient *ProductClient
}

func (s *Service) AddToCart(userID uint, productID uint, qty int, authorizationHeader string) error {
	if productID == 0 {
		return errors.New("product id is required")
	}

	if qty < 1 {
		return errors.New("quantity must be at least 1")
	}

	product, err := s.ProductClient.GetProduct(productID, authorizationHeader)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Stock < qty {
		return errors.New("insufficient stock")
	}

	var cart Cart

	err = s.DB.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("failed to fetch cart")
		}

		cart = Cart{UserID: userID}
		if err := s.DB.Create(&cart).Error; err != nil {
			return errors.New("failed to create cart")
		}
	}

	var item CartItem

	err = s.DB.Where("cart_id = ? AND product_id = ?", cart.ID, productID).
		First(&item).Error

	if err == nil {

		newQty := item.Quantity + qty

		if product.Stock < newQty {
			return errors.New("insufficient stock")
		}

		item.Quantity = newQty
		return s.DB.Save(&item).Error
	}

	item = CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  qty,
	}

	return s.DB.Create(&item).Error
}

func (s *Service) ReduceItem(userID uint, productID uint) error {
	if productID == 0 {
		return errors.New("product id is required")
	}

	var cart Cart
	if err := s.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return errors.New("cart not found")
	}

	var item CartItem

	if err := s.DB.Where("cart_id = ? AND product_id = ?", cart.ID, productID).
		First(&item).Error; err != nil {
		return errors.New("item not found")
	}

	if item.Quantity > 1 {
		item.Quantity--
		return s.DB.Save(&item).Error
	}

	return s.DB.Delete(&item).Error
}

func (s *Service) RemoveItem(userID uint, productID uint) error {
	if productID == 0 {
		return errors.New("product id is required")
	}

	var cart Cart
	if err := s.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return errors.New("cart not found")
	}

	result := s.DB.Where("cart_id = ? AND product_id = ?", cart.ID, productID).
		Delete(&CartItem{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}

func (s *Service) ClearCart(userID uint) error {

	var cart Cart
	if err := s.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return errors.New("failed to fetch cart")
	}

	return s.DB.Where("cart_id = ?", cart.ID).Delete(&CartItem{}).Error
}

func (s *Service) GetCart(userID uint, authorizationHeader string) (*CartResponse, error) {

	var cart Cart
	if err := s.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return emptyCartResponse(), nil
		}

		return nil, errors.New("failed to fetch cart")
	}

	var items []CartItem
	if err := s.DB.Where("cart_id = ?", cart.ID).Find(&items).Error; err != nil {
		return nil, errors.New("failed to fetch cart items")
	}

	responseItems := []CartItemResponse{}
	var subtotal float64

	for _, item := range items {

		product, err := s.ProductClient.GetProduct(item.ProductID, authorizationHeader)
		if err != nil {
			continue
		}

		image := ""
		if len(product.ImageURLs) > 0 {
			image = product.ImageURLs[0]
		}

		total := product.Price * float64(item.Quantity)
		subtotal += total

		responseItems = append(responseItems, CartItemResponse{
			ProductID:   product.ID,
			Title:       product.Title,
			Description: product.Description,
			ImageURL:    image,
			Price:       product.Price,
			Quantity:    item.Quantity,
			Total:       total,

			ExpectedDelivery:  utils.GetExpectedDelivery(),
		})
	}
	delivery := ""

if len(responseItems) > 0 {
	delivery = utils.GetExpectedDelivery()
}

	return &CartResponse{
		Items:            responseItems,
		Subtotal:         subtotal,
		Total:            subtotal,
		ExpectedDelivery:  delivery,
	}, nil
}

func emptyCartResponse() *CartResponse {
	return &CartResponse{
		Items:            []CartItemResponse{},
		Subtotal:         0,
		Total:            0,
		ExpectedDelivery: "",
	}
}

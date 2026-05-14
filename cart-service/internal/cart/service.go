package cart

import (
	"errors"
	"log"
	"cart-service/pkg/utils"
	"gorm.io/gorm"
)

type Service struct {
	Repo          *Repository
	ProductClient *ProductClient
}

func NewService(
	repo *Repository,
	productClient *ProductClient,
) *Service {
	return &Service{
		Repo:          repo,
		ProductClient: productClient,
	}
}

func (s *Service) AddToCart(
	userID uint,
	productID uint,
	qty int,
	authorizationHeader string,
) (*CartResponse, error) {



	product, err := s.ProductClient.GetProduct(
		productID,
		authorizationHeader,
	)

	if err != nil {
		return nil, errors.New("product not found")
	}

	if product.Stock < qty {
		return nil, errors.New("insufficient stock")
	}

	cart, err := s.Repo.GetCartByUserID(userID)

	if err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("failed to fetch cart")
		}

		cart, err = s.Repo.CreateCart(userID)

		if err != nil {
			return nil, errors.New("failed to create cart")
		}
	}

	item, err := s.Repo.GetCartItem(cart.ID, productID)

	if err == nil {

		newQty := item.Quantity + qty

		if newQty > product.Stock {
			return nil, errors.New("insufficient stock")
		}

		item.Quantity = newQty

		if err := s.Repo.UpdateCartItem(item); err != nil {
			return nil, err
		}

		return s.GetCart(userID, authorizationHeader)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	newItem := &CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  qty,
	}

	if err := s.Repo.CreateCartItem(newItem); err != nil {
		return nil, err
	}

	return s.GetCart(userID, authorizationHeader)
}

func (s *Service) ReduceItem(userID uint,productID uint) error {
	cart, err := s.Repo.GetCartByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	item, err := s.Repo.GetCartItem(cart.ID, productID)

	if err != nil {
		return errors.New("item not found")
	}

	if item.Quantity > 1 {
		item.Quantity--
		return s.Repo.UpdateCartItem(item)
	}

	return s.Repo.DeleteCartItem(item)
}

func (s *Service) RemoveItem(
	userID uint,
	productID uint,
) error {

	if productID == 0 {
		return errors.New("product id is required")
	}

	cart, err := s.Repo.GetCartByUserID(userID)

	if err != nil {
		return errors.New("cart not found")
	}

	return s.Repo.DeleteCartItemByProductID(
		cart.ID,
		productID,
	)
}

func (s *Service) ClearCart(userID uint) error {

	cart, err := s.Repo.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return errors.New("failed to fetch cart")
	}

	return s.Repo.ClearCart(cart.ID)
}

func (s *Service) GetCart(userID uint, authorizationHeader string) (*CartResponse, error) {

	cart, err := s.Repo.GetCartByUserID(userID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return emptyCartResponse(), nil
		}

		return nil, errors.New("failed to fetch cart")
	}

	items, err := s.Repo.GetCartItems(cart.ID)
	if err != nil {
		return nil, errors.New("failed to fetch cart items")
	}

	responseItems := []CartItemResponse{}
	var subtotal float64

	for _, item := range items {

		product, err := s.ProductClient.GetProduct(item.ProductID, authorizationHeader)
		if err != nil {
    log.Printf("failed to fetch product %d: %v", item.ProductID, err)
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

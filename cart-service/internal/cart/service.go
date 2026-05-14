package cart

import (
	"errors"

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

	products, err := s.ProductClient.GetProducts(
		[]uint{productID},
		authorizationHeader,
	)
	if err != nil {
		return nil, errors.New("product not found")
	}

	product, exists := products[productID]
	if !exists {
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
	} else {
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
	}

	items, err := s.Repo.GetCartItems(cart.ID)
	if err != nil {
		return nil, err
	}

	return s.buildCartResponse(items, authorizationHeader)
}

func (s *Service) ReduceItem(
	userID uint,
	productID uint,
	authorizationHeader string,
) (*CartItemResponse, error) {

	cart, err := s.Repo.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	item, err := s.Repo.GetCartItem(cart.ID, productID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if item.Quantity > 1 {
		item.Quantity--
		if err := s.Repo.UpdateCartItem(item); err != nil {
			return nil, err
		}

		products, err := s.ProductClient.GetProducts(
			[]uint{productID},
			authorizationHeader,
		)
		if err != nil {
			return nil, errors.New("product not found")
		}

		product, exists := products[productID]
		if !exists {
			return nil, errors.New("product not found")
		}

		resp := cartItemResponseFromProduct(product, item.Quantity)
		return &resp, nil
	}

	if err := s.Repo.DeleteCartItem(item); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Service) IncreaseItem(
	userID uint,
	productID uint,
	authorizationHeader string,
) (*CartItemResponse, error) {

	cart, err := s.Repo.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("cart not found")
	}

	item, err := s.Repo.GetCartItem(cart.ID, productID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	products, err := s.ProductClient.GetProducts(
		[]uint{productID},
		authorizationHeader,
	)
	if err != nil {
		return nil, errors.New("product not found")
	}

	product, exists := products[productID]
	if !exists {
		return nil, errors.New("product not found")
	}

	if item.Quantity+1 > product.Stock {
		return nil, errors.New("insufficient stock")
	}

	if err := s.Repo.IncrementCartItemQty(cart.ID, productID); err != nil {
		return nil, err
	}

	updatedItem, err := s.Repo.GetCartItem(cart.ID, productID)
	if err != nil {
		return nil, errors.New("failed to fetch updated item")
	}

	resp := cartItemResponseFromProduct(product, updatedItem.Quantity)
	return &resp, nil
}

func (s *Service) RemoveItem(
	userID uint,
	productID uint,
) error {

	cart, err := s.Repo.GetCartByUserID(userID)
	if err != nil {
		return errors.New("cart not found")
	}

	return s.Repo.DeleteCartItemByProductID(cart.ID, productID)
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

func (s *Service) GetCart(
	userID uint,
	authorizationHeader string,
) (*CartResponse, error) {

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

	return s.buildCartResponse(items, authorizationHeader)
}

func cartItemResponseFromProduct(product ProductDTO, quantity int) CartItemResponse {
	image := ""
	if len(product.ImageURLs) > 0 {
		image = product.ImageURLs[0]
	}

	total := product.Price * float64(quantity)

	return CartItemResponse{
		ProductID:        product.ID,
		Title:            product.Title,
		Description:      product.Description,
		ImageURL:         image,
		Price:            product.Price,
		Quantity:         quantity,
		Total:            total,
		ExpectedDelivery: utils.GetExpectedDelivery(),
	}
}

func (s *Service) buildCartResponse(
	items []CartItem,
	authorizationHeader string,
) (*CartResponse, error) {

	if len(items) == 0 {
		return emptyCartResponse(), nil
	}

	productIDs := make([]uint, len(items))
	for i, item := range items {
		productIDs[i] = item.ProductID
	}

	productsMap, err := s.ProductClient.GetProducts(productIDs, authorizationHeader)
	if err != nil {
		return nil, err
	}

	responseItems := []CartItemResponse{}
	var subtotal float64

	for _, item := range items {
		product, exists := productsMap[item.ProductID]
		if !exists {
			continue
		}

		line := cartItemResponseFromProduct(product, item.Quantity)
		subtotal += line.Total
		responseItems = append(responseItems, line)
	}

	delivery := ""
	if len(responseItems) > 0 {
		delivery = utils.GetExpectedDelivery()
	}

	return &CartResponse{
		Items:            responseItems,
		Subtotal:         subtotal,
		Total:            subtotal,
		ExpectedDelivery: delivery,
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
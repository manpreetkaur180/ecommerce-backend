package order

import "errors"

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) CreateOrder(
	userID uint,
	token string,
	req CreateOrderRequest,
) (*Order, error) {

	cartItems, err := FetchCartItems(token)
	if err != nil {
		return nil, errors.New("failed to fetch cart")
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	addresses, err := FetchUserAddresses(token)
	if err != nil {
		return nil, errors.New("failed to fetch addresses")
	}

	var selectedAddress *AddressResponse

	for _, addr := range addresses {
		if addr.ID == req.AddressID {
			selectedAddress = &addr
			break
		}
	}

	if selectedAddress == nil {
		return nil, errors.New("address not found")
	}

	for _, item := range cartItems {

		inventory, err := FetchInventory(item.ProductID)
		if err != nil {
			return nil, errors.New("failed to validate inventory")
		}

		if inventory.Stock < item.Quantity {
			return nil, errors.New(
				item.ProductName + " is out of stock",
			)
		}
	}

	tx := s.Repo.DB.Begin()

	txRepo := s.Repo.WithTx(tx)

	order := Order{
		UserID: userID,

		AddressID: selectedAddress.ID,

		FullName: selectedAddress.FullName,
		Phone:    selectedAddress.Phone,

		AddressLine1: selectedAddress.AddressLine1,
		AddressLine2: selectedAddress.AddressLine2,
		Landmark:     selectedAddress.Landmark,
		City:         selectedAddress.City,
		State:        selectedAddress.State,
		Country:      selectedAddress.Country,
		Pincode:      selectedAddress.Pincode,

		Status: OrderPending,
	}

	var total float64

	for _, item := range cartItems {

		subtotal := float64(item.Quantity) * item.Price

		order.Items = append(order.Items, OrderItem{
			ProductID: item.ProductID,
			SellerID:  item.SellerID,

			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,

			Quantity: item.Quantity,
			Price:    item.Price,
			Subtotal: subtotal,
		})

		total += subtotal
	}

	order.TotalAmount = total

	if err := txRepo.CreateOrder(&order); err != nil {
		tx.Rollback()
		return nil, errors.New("failed to create order")
	}

	for _, item := range cartItems {

		if err := ReduceInventory(
			item.ProductID,
			item.Quantity,
		); err != nil {

			tx.Rollback()

			return nil, errors.New(
				"failed to update inventory",
			)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("failed to commit order")
	}

	return &order, nil
}

func (s *Service) GetUserOrders(userID uint) ([]Order, error) {
	return s.Repo.FindOrdersByUser(userID)
}

func (s *Service) GetOrderByID(id uint) (*Order, error) {
	return s.Repo.FindOrderByID(id)
}

func (s *Service) ConfirmOrder(id uint) error {

	order, err := s.Repo.FindOrderByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != OrderPending {
		return errors.New("only pending orders can be confirmed")
	}

	order.Status = OrderConfirmed

	return s.Repo.SaveOrder(order)
}

func (s *Service) RejectOrder(id uint) error {

	order, err := s.Repo.FindOrderByID(id)
	if err != nil {
		return errors.New("order not found")
	}

	if order.Status != OrderPending {
		return errors.New("only pending orders can be rejected")
	}

	order.Status = OrderRejected

	return s.Repo.SaveOrder(order)
}

func (s *Service) CancelOrder(
	userID uint,
	orderID uint,
) error {

	order, err := s.Repo.FindOrderByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	if order.UserID != userID {
		return errors.New("unauthorized")
	}

	if order.Status == OrderDelivered {
		return errors.New("delivered order cannot be cancelled")
	}

	order.Status = OrderCancelled

	return s.Repo.SaveOrder(order)
}
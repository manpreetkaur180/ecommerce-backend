package order

import (
    "ecommerce-backend/internal/cart"
)

type Service struct {
    repo     *Repository
    cartRepo *cart.Repository
}

func NewService(repo *Repository, cartRepo *cart.Repository) *Service {
    return &Service{repo, cartRepo}
}

func (s *Service) CreateOrder(userID uint) error {
    cartData, err := s.cartRepo.GetCartByUser(userID)
    if err != nil {
        return err
    }

    var total float64
    var items []OrderItem

    for _, item := range cartData.Items {
        price := 100.0 // TEMP (replace later with product price)

        total += price * float64(item.Quantity)

        items = append(items, OrderItem{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            Price:     price,
        })
    }

    order := &Order{
        UserID:      userID,
        TotalAmount: total,
        Status:      "created",
        Items:       items,
    }

    return s.repo.CreateOrder(order)
}
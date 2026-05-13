package cart

import "time"

type Cart struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItem struct {
	ID        uint `gorm:"primaryKey"`
	CartID    uint `gorm:"index;not null"`
	ProductID uint `gorm:"index;not null"`
	Quantity  int  `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type AddToCartRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type ReduceCartRequest struct {
	ProductID uint `json:"product_id"`
}

type CartItemResponse struct {
	ProductID        uint    `json:"product_id"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	ImageURL         string  `json:"image_url"`
	Price            float64 `json:"price"`
	Quantity         int     `json:"quantity"`
	Total            float64 `json:"total"`
	ExpectedDelivery string  `json:"expected_delivery"`
}
type CartResponse struct {
	Items            []CartItemResponse `json:"items"`
	Subtotal         float64            `json:"subtotal"`
	Total            float64            `json:"total"`
	ExpectedDelivery string             `json:"expected_delivery"`
}
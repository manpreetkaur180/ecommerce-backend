package cart


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
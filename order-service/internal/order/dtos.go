package order

type CreateOrderRequest struct {
	AddressID uint `json:"address_id" validate:"required"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required"`
}
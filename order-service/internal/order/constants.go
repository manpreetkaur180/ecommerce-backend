package order

type OrderStatus string

const (
	OrderPending        OrderStatus = "PENDING"
	OrderConfirmed      OrderStatus = "CONFIRMED"
	OrderRejected       OrderStatus = "REJECTED"
	OrderAgentAssigned  OrderStatus = "AGENT_ASSIGNED"
	OrderPacked         OrderStatus = "PACKED"
	OrderOutForDelivery OrderStatus = "OUT_FOR_DELIVERY"
	OrderDelivered      OrderStatus = "DELIVERED"
	OrderCancelled      OrderStatus = "CANCELLED"
)
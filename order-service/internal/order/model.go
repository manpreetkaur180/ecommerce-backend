package order

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	UserID uint `gorm:"index"`

	AddressID uint

	FullName string
	Phone    string

	AddressLine1 string
	AddressLine2 string
	Landmark     string
	City         string
	State        string
	Country      string
	Pincode      string

	Status OrderStatus `gorm:"type:varchar(40);default:'PENDING'"`

	TotalAmount float64

	Items []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model

	OrderID uint `gorm:"index"`

	ProductID uint
	SellerID  uint

	ProductName string
	ProductImage string

	Quantity int
	Price    float64
	Subtotal float64
}
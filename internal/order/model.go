package order

import "time"

type Order struct {
    ID          uint        `gorm:"primaryKey"`
    UserID      uint
    TotalAmount float64
    Status      string
    Items       []OrderItem `gorm:"foreignKey:OrderID"`
    CreatedAt   time.Time
}

type OrderItem struct {
    ID        uint `gorm:"primaryKey"`
    OrderID   uint
    ProductID uint
    Quantity  int
    Price     float64
}
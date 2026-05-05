package cart

import "time"

type Cart struct {
    ID        uint       `gorm:"primaryKey"`
    UserID    uint       `gorm:"unique"` // one cart per user
    Items     []CartItem `gorm:"foreignKey:CartID"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type CartItem struct {
    ID        uint `gorm:"primaryKey"`
    CartID    uint
    ProductID uint
    Quantity  int
}
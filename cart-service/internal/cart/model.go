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

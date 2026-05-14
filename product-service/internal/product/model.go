package product

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID uint `gorm:"primaryKey"`

	SellerID uint `gorm:"not null;index"`

	Title string `gorm:"not null"`

	Description string `gorm:"type:text"`

	Price float64 `gorm:"not null"`

	Stock int `gorm:"default:0"`

	Category string

	ImageURLs datatypes.JSONSlice[string] `gorm:"type:json"`

	Offers string `gorm:"type:text"`

	Rating float64 `gorm:"default:0"`

	ReturnAvailable bool `gorm:"default:false"`

	Warranty string `gorm:"type:text"`

	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

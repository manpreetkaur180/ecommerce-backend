package product

import "time"

type Product struct {
	ID uint `gorm:"primaryKey"`

	SellerID uint `gorm:"not null;index"`

	Title string `gorm:"not null"`

	Description string `gorm:"type:text"`

	Price float64 `gorm:"not null"`

	Stock int `gorm:"default:0"`

	Category string

	ImageURL string

	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateProductRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

type UpdateProductRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}
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

	ImageURLs datatypes.JSONSlice[string] `gorm:"type:json" json:"image_urls"`

	Offers string  `gorm:"type:text" json:"offers"`
	Rating float64 `gorm:"default:0" json:"rating"`

	ReturnAvailable bool `gorm:"default:false" json:"return_available"`

	Warranty string `gorm:"type:text" json:"warranty,omitempty"`

	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateProductRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	Category    string   `json:"category"`
	ImageURLs   []string `json:"image_urls"`
	Offers      string   `json:"offers"`
	Rating      float64  `json:"rating"`

	ReturnAvailable bool `json:"return_available"`

	Warranty string `json:"warranty"`
}

type UpdateProductRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       *int     `json:"stock"`
	Category    string   `json:"category"`
	ImageURLs   []string `json:"image_urls"`
	Offers      string   `json:"offers"`
	IsActive    *bool    `json:"is_active"`
	Rating      *float64 `json:"rating"`

	ReturnAvailable *bool `json:"return_available"`

	Warranty string `json:"warranty"`
}

type BuyerProductResponse struct {
	ID               uint    `json:"id"`
	Title            string  `json:"title"`
	ImageURL         string  `json:"image_url"`
	Description      string  `json:"description"`
	Price            float64 `json:"price"`
	Offers           string  `json:"offers"`
	ExpectedDelivery string  `json:"expected_delivery"`
}
type BuyerProductDetailResponse struct {
	ID                 uint     `json:"id"`

	Title              string   `json:"title"`

	Description        string   `json:"description"`

	Price              float64  `json:"price"`

	Category           string   `json:"category"`

	ImageURLs          []string `json:"image_urls"`

	Offers             string   `json:"offers"`

	Rating             float64  `json:"rating"`

	ReturnAvailable    bool     `json:"return_available"`

	Warranty           string   `json:"warranty,omitempty"`

	InStock            bool     `json:"in_stock"`

	AvailableQuantity  int      `json:"available_quantity"`

	ExpectedDelivery   string   `json:"expected_delivery"`
}
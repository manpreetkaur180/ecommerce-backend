package product

import (
	"time"

	"gorm.io/datatypes"
)

type Product struct {
	ID uint `gorm:"primaryKey" json:"id"`

	SellerID uint `gorm:"not null;index" json:"seller_id"`

	Title string `gorm:"not null" json:"title"`

	Description string `gorm:"type:text" json:"description"`

	Price float64 `gorm:"not null" json:"price"`

	Stock int `gorm:"default:0" json:"stock"`

	Category string `json:"category"`

	ImageURLs datatypes.JSONSlice[string] `gorm:"type:json" json:"image_urls"`

	Offers string `gorm:"type:text" json:"offers"`

	Rating float64 `gorm:"default:0" json:"rating"`

	ReturnAvailable bool `gorm:"default:false" json:"return_available"`

	Warranty string `gorm:"type:text" json:"warranty,omitempty"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Stock       int      `json:"stock"`
	Category    string   `json:"category"`
	ImageURLs   []string `json:"image_urls"`
	Offers      string   `json:"offers"`

	ReturnAvailable bool `json:"return_available"`

	Warranty string `json:"warranty"`
}

type UpdateProductRequest struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Price       *float64  `json:"price"`
	Stock       *int      `json:"stock"`
	Category    *string   `json:"category"`
	ImageURLs   *[]string `json:"image_urls"`
	Offers      *string   `json:"offers"`
	IsActive    *bool     `json:"is_active"`

	ReturnAvailable *bool `json:"return_available"`

	Warranty *string `json:"warranty"`
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
	ID uint `json:"id"`

	Title string `json:"title"`

	Description string `json:"description"`

	Price float64 `json:"price"`

	Category string `json:"category"`

	ImageURLs []string `json:"image_urls"`

	Offers string `json:"offers"`

	Rating float64 `json:"rating"`

	ReturnAvailable bool `json:"return_available"`

	Warranty string `json:"warranty,omitempty"`

	InStock bool `json:"in_stock"`

	Stock int `json:"stock"`

	ExpectedDelivery string `json:"expected_delivery"`
}
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type BuyerProductsPaginatedResponse struct {
	Products   []BuyerProductResponse `json:"products"`
	Pagination PaginationMeta         `json:"pagination"`
}

type SellerProductsPaginatedResponse struct {
	Products   []Product      `json:"products"`
	Pagination PaginationMeta `json:"pagination"`
}

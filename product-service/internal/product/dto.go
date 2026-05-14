package product

import "time"

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

type BulkProductIDsRequest struct {
	ProductIDs []uint `json:"product_ids"`
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

type ProductResponse struct {
	ID              uint      `json:"id"`
	SellerID        uint      `json:"seller_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	Stock           int       `json:"stock"`
	Category        string    `json:"category"`
	ImageURLs       []string  `json:"image_urls"`
	Offers          string    `json:"offers"`
	Rating          float64   `json:"rating"`
	ReturnAvailable bool      `json:"return_available"`
	Warranty        string    `json:"warranty,omitempty"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type SellerProductsPaginatedResponse struct {
	Products   []ProductResponse `json:"products"`
	Pagination PaginationMeta    `json:"pagination"`
}

func productToResponse(p Product) ProductResponse {
	return ProductResponse{
		ID:              p.ID,
		SellerID:        p.SellerID,
		Title:           p.Title,
		Description:     p.Description,
		Price:           p.Price,
		Stock:           p.Stock,
		Category:        p.Category,
		ImageURLs:       []string(p.ImageURLs),
		Offers:          p.Offers,
		Rating:          p.Rating,
		ReturnAvailable: p.ReturnAvailable,
		Warranty:        p.Warranty,
		IsActive:        p.IsActive,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

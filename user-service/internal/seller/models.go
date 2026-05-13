package seller

import "time"

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type SellerApplication struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null;index"`

	BusinessName string `gorm:"not null"`

	BusinessDescription string

	GSTIN string `gorm:"not null"`

	AadharNumber string `gorm:"not null"`

	Status string `gorm:"default:'pending'"`

	AdminNote string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ApplySellerRequest struct {
	BusinessName        string `json:"business_name"`
	BusinessDescription string `json:"business_description"`
	GSTIN               string `json:"gstin"`
	AadharNumber        string `json:"aadhar_number"`
}

type SellerApplicationListResponse struct {
	ID                  uint      `json:"id"`
	UserID              uint      `json:"user_id"`
	BusinessName        string    `json:"business_name"`
	BusinessDescription string    `json:"business_description"`
	GSTIN               string    `json:"gstin"`
	AadharNumber        string    `json:"aadhar_number"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type SellerApplicationDecisionRequest struct {
	AdminNote string `json:"admin_note"`
}

type SellerApplicationDecisionResponse struct {
	ID                  uint      `json:"id"`
	UserID              uint      `json:"user_id"`
	BusinessName        string    `json:"business_name"`
	BusinessDescription string    `json:"business_description"`
	GSTIN               string    `json:"gstin"`
	AadharNumber        string    `json:"aadhar_number"`
	Status              string    `json:"status"`
	AdminNote           string    `json:"admin_note"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

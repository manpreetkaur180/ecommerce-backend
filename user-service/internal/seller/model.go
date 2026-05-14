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

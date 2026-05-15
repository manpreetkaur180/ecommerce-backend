package user

import "time"

const (
	RoleBuyer  = "buyer"
	RoleSeller = "seller"
	RoleAdmin  = "admin"
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Email      string `gorm:"uniqueIndex;not null"`
	Phone      string `gorm:"uniqueIndex;not null"`
	Password   string `gorm:"not null"`
	IsVerified bool   `gorm:"default:false"`

	Role string `gorm:"type:varchar(20);default:'buyer'"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type EmailVerification struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
}

type PasswordResetVerification struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type PasswordReset struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type PasswordUpdate struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
type Address struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null;index"`

	FullName string `gorm:"not null"`

	Phone string `gorm:"not null"`

	AddressLine1 string `gorm:"not null"`

	AddressLine2 string

	Landmark string

	City string `gorm:"not null"`

	State string `gorm:"not null"`

	Country string `gorm:"not null"`

	Pincode string `gorm:"not null"`

	AddressType string `gorm:"default:'home'"` 

	IsDefault bool `gorm:"default:false"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

package user

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"not null"`
	Email      string    `gorm:"uniqueIndex;not null"`
	Phone      string    `gorm:"uniqueIndex;not null"`
	Password   string    `gorm:"not null"`

	CreatedAt  time.Time
	UpdatedAt  time.Time
}


type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
type OTPLoginRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}
type OTPData struct {
	Code    string
	Expires time.Time
}
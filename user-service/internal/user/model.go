package user

import "time"

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

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UpdatePasswordRequest struct {
	Email string `json:"email"`
}

type UpdatePasswordFormRequest struct {
	Token              string `json:"token"`
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

type PasswordUpdate struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	TokenHash string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

type UpdatePasswordConfirmRequest struct {
	Token              string `json:"token"`
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	ConfirmNewPassword string `json:"confirm_new_password"`
}

const (
	RoleBuyer  = "buyer"
	RoleSeller = "seller"
	RoleAdmin  = "admin"
)

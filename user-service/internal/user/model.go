package user

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"not null"`
	Email      string    `gorm:"uniqueIndex;not null"`
	Phone      string    `gorm:"uniqueIndex"`
	Password   string    `gorm:"not null"`
	IsVerified bool      `gorm:"default:false"`
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
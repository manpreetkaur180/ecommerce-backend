package user

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {

	adminEmail := os.Getenv("ADMIN_EMAIL")

	var existing User

	err := db.Where(
		"email = ?",
		adminEmail,
	).First(&existing).Error

	// already exists
	if err == nil {
		log.Println("Admin already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(os.Getenv("ADMIN_PASSWORD")),
		10,
	)

	if err != nil {
		log.Println("Failed to hash admin password")
		return
	}

	admin := User{
		Name:       os.Getenv("ADMIN_NAME"),
		Email:      os.Getenv("ADMIN_EMAIL"),
		Phone:      os.Getenv("ADMIN_PHONE"),
		Password:   string(hashedPassword),
		IsVerified: true,
		Role:       RoleAdmin,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Println("Failed to seed admin:", err)
		return
	}

	log.Println("Admin seeded successfully")
}

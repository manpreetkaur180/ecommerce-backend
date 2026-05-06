package user

import (
	"errors"

	"user-service/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) Register(req RegisterRequest) (*User, error) {

	// 1. basic validation
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := utils.ValidatePhone(req.Phone); err != nil {
		return nil, err
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	if req.Password != req.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}

	// 2. check email exists
	var existing User
	err := s.DB.Where("email = ?", req.Email).First(&existing).Error
	if err == nil {
		return nil, errors.New("email already exists")
	}

	// 3. hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 4. create user
	user := User{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashed),
	}

	err = s.DB.Create(&user).Error
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return &user, nil
}
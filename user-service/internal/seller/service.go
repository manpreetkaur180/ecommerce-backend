package seller

import (
	"errors"
	"user-service/internal/user"

	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}
func (s *Service) ApplySeller(
	userID uint,
	req ApplySellerRequest,
) error {

	// check user exists
	var existingUser user.User

	if err := s.DB.First(&existingUser, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// already seller
	if existingUser.IsSeller {
		return errors.New("user is already a seller")
	}

	// already applied
	var existingApplication SellerApplication

	err := s.DB.Where(
		"user_id = ? AND status = ?",
		userID,
		StatusPending,
	).First(&existingApplication).Error

	if err == nil {
		return errors.New("seller application already pending")
	}

	// validation
	if req.BusinessName == "" {
		return errors.New("business name is required")
	}

	if req.GSTIN == "" {
		return errors.New("gstin is required")
	}

	if req.AadharNumber == "" {
		return errors.New("aadhar number is required")
	}

	// create application
	application := SellerApplication{
		UserID:               userID,
		BusinessName:         req.BusinessName,
		BusinessDescription:  req.BusinessDescription,
		GSTIN:                req.GSTIN,
		AadharNumber:         req.AadharNumber,
		Status:               StatusPending,
	}

	if err := s.DB.Create(&application).Error; err != nil {
		return errors.New("failed to create seller application")
	}

	return nil
}
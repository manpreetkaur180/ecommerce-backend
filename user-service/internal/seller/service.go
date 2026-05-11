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
func (s *Service) GetAllApplications() ([]SellerApplication, error) {

	var applications []SellerApplication

	if err := s.DB.Order("created_at desc").
		Find(&applications).Error; err != nil {

		return nil, errors.New("failed to fetch applications")
	}

	return applications, nil
}
func (s *Service) ApproveApplication(
	applicationID uint,
) error {

	var application SellerApplication

	if err := s.DB.First(
		&application,
		applicationID,
	).Error; err != nil {

		return errors.New("application not found")
	}

	if application.Status != StatusPending {
		return errors.New("application already processed")
	}

	// fetch user
	var existingUser user.User

	if err := s.DB.First(
		&existingUser,
		application.UserID,
	).Error; err != nil {

		return errors.New("user not found")
	}

	// update seller status
	existingUser.IsSeller = true

	if err := s.DB.Save(&existingUser).Error; err != nil {
		return errors.New("failed to update user seller status")
	}

	// approve application
	application.Status = StatusApproved

	if err := s.DB.Save(&application).Error; err != nil {
		return errors.New("failed to approve application")
	}

	return nil
}
func (s *Service) RejectApplication(
	applicationID uint,
	adminNote string,
) error {

	var application SellerApplication

	if err := s.DB.First(
		&application,
		applicationID,
	).Error; err != nil {

		return errors.New("application not found")
	}

	if application.Status != StatusPending {
		return errors.New("application already processed")
	}

	application.Status = StatusRejected
	application.AdminNote = adminNote

	if err := s.DB.Save(&application).Error; err != nil {
		return errors.New("failed to reject application")
	}

	return nil
}
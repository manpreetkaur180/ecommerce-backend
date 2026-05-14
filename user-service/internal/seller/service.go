package seller

import (
	"errors"
	"strings"
	"user-service/internal/user"
)

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) ApplySeller(
	userID uint,
	req ApplySellerRequest,
) error {

	existingUser, err := s.Repo.FirstUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if existingUser.Role == user.RoleSeller {
		return errors.New("user is already a seller")
	}

	if _, err := s.Repo.FirstPendingApplicationByUserID(userID); err == nil {
		return errors.New("seller application already pending")
	}

	if req.BusinessName == "" {
		return errors.New("business name is required")
	}

	if req.GSTIN == "" {
		return errors.New("gstin is required")
	}

	if req.AadharNumber == "" {
		return errors.New("aadhar number is required")
	}

	application := SellerApplication{
		UserID:              userID,
		BusinessName:        req.BusinessName,
		BusinessDescription: req.BusinessDescription,
		GSTIN:               req.GSTIN,
		AadharNumber:        req.AadharNumber,
		Status:              StatusPending,
	}

	if err := s.Repo.CreateSellerApplication(&application); err != nil {
		return errors.New("failed to create seller application")
	}

	return nil
}

func (s *Service) GetAllApplications() ([]SellerApplication, error) {
	applications, err := s.Repo.FindAllSellerApplicationsOrdered()
	if err != nil {
		return nil, errors.New("failed to fetch applications")
	}
	return applications, nil
}

func (s *Service) ApproveApplication(
	applicationID uint,
	adminNote string,
) (*SellerApplication, error) {

	adminNote = strings.TrimSpace(adminNote)
	if adminNote == "" {
		return nil, errors.New("admin note is required")
	}

	return s.Repo.ApproveApplicationTx(applicationID, adminNote)
}

func (s *Service) RejectApplication(
	applicationID uint,
	adminNote string,
) (*SellerApplication, error) {

	adminNote = strings.TrimSpace(adminNote)
	if adminNote == "" {
		return nil, errors.New("admin note is required")
	}

	application, err := s.Repo.FirstSellerApplicationByID(applicationID)
	if err != nil {
		return nil, errors.New("application not found")
	}

	if application.Status != StatusPending {
		return nil, errors.New("application already processed")
	}

	application.Status = StatusRejected
	application.AdminNote = adminNote

	if err := s.Repo.SaveSellerApplication(application); err != nil {
		return nil, errors.New("failed to reject application")
	}

	return application, nil
}

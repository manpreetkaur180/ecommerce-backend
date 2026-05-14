package seller

import (
	"errors"
	"user-service/internal/user"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) FirstUserByID(id uint) (*user.User, error) {
	var u user.User
	if err := r.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FirstPendingApplicationByUserID(userID uint) (*SellerApplication, error) {
	var app SellerApplication
	err := r.DB.Where(
		"user_id = ? AND status = ?",
		userID,
		StatusPending,
	).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *Repository) CreateSellerApplication(a *SellerApplication) error {
	return r.DB.Create(a).Error
}

func (r *Repository) FindAllSellerApplicationsOrdered() ([]SellerApplication, error) {
	var applications []SellerApplication
	if err := r.DB.Order("created_at desc").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *Repository) ApproveApplicationTx(
	applicationID uint,
	adminNote string,
) (*SellerApplication, error) {

	var approvedApplication SellerApplication

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var application SellerApplication

		if err := tx.First(&application, applicationID).Error; err != nil {
			return errors.New("application not found")
		}

		if application.Status != StatusPending {
			return errors.New("application already processed")
		}

		var existingUser user.User

		if err := tx.First(&existingUser, application.UserID).Error; err != nil {
			return errors.New("user not found")
		}

		existingUser.Role = user.RoleSeller

		if err := tx.Save(&existingUser).Error; err != nil {
			return errors.New("failed to update user role")
		}

		application.Status = StatusApproved
		application.AdminNote = adminNote

		if err := tx.Save(&application).Error; err != nil {
			return errors.New("failed to approve application")
		}

		approvedApplication = application
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &approvedApplication, nil
}

func (r *Repository) FirstSellerApplicationByID(id uint) (*SellerApplication, error) {
	var application SellerApplication
	if err := r.DB.First(&application, id).Error; err != nil {
		return nil, err
	}
	return &application, nil
}

func (r *Repository) SaveSellerApplication(a *SellerApplication) error {
	return r.DB.Save(a).Error
}

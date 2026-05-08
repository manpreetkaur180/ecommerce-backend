package services

import "notification-service/internal/clients"

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) HandleUserRegistered(
	name string,
	email string,
) error {

	return clients.SendWelcomeEmail(
		name,
		email,
	)
}

func (s *NotificationService) HandleOTPEmail(name, email, otp string) error {
	return clients.SendOTPEmail(
		name,
		email,
		otp,
	)
}

func (s *NotificationService) HandleVerificationEmail(name, email, link string) error {
	return clients.SendVerificationEmail(
		name,
		email,
		link,
	)
}

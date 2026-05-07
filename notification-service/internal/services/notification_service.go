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
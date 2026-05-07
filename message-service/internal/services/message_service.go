package services

import (
	"fmt"

	"message-service/internal/models"
	"message-service/internal/templates"
)

type Sender interface {
	Send(message models.Message) error
}

type MessageService struct {
	sender Sender
}

func NewMessageService(sender Sender) *MessageService {
	return &MessageService{
		sender: sender,
	}
}

func (s *MessageService) SendEmail(
	req models.EmailRequest,
) error {

	content := ""

	switch req.Template {

	case "otp":

		name := fmt.Sprintf("%v", req.Data["name"])
		otp := fmt.Sprintf("%v", req.Data["otp"])

		content = templates.OTPTemplate(
			name,
			otp,
		)

	default:
		content = "No template found"
	}

	message := models.Message{
		To:       req.To,
		Subject:  req.Subject,
		Content:  content,
		Template: req.Template,
	}

	return s.sender.Send(message)
}
package services

import (
	"fmt"

	"message-service/internal/models"
	"message-service/internal/templates"
	"errors"
)

type Sender interface {
	Send(message models.Message) error
}

type MessageService struct {
	emailSender Sender
	smsSender   Sender
}

func NewMessageService(
	emailSender Sender,
	smsSender Sender,
) *MessageService {
	return &MessageService{
		emailSender: emailSender,
		smsSender:   smsSender,
	}
}

func (s *MessageService) SendEmail(req models.EmailRequest) error {

	var content string

	switch req.Template {

	case "otp":

		content = templates.OTPTemplate(
			req.Data["name"].(string),
			req.Data["otp"].(string),
		)

	case "welcome":

		content = templates.WelcomeTemplate(
			req.Data["name"].(string),
		)

	case "verify_email":

		content = templates.VerifyEmailTemplate(
			req.Data["name"].(string),
			req.Data["link"].(string),
		)

	case "forgot_password_verify":

		content = templates.ForgotPasswordVerifyTemplate(
			req.Data["name"].(string),
			req.Data["link"].(string),
		)

	case "reset_password":

		content = templates.ResetPasswordTemplate(
			req.Data["name"].(string),
			req.Data["link"].(string),
		)

	case "update_password":

		content = templates.UpdatePasswordTemplate(
			req.Data["name"].(string),
			req.Data["link"].(string),
		)

	default:

		return errors.New("no template found")
	}
	message := models.Message{
		To:       req.To,
		Subject:  req.Subject,
		Content:  content,
		Template: req.Template,
	}

	return s.emailSender.Send(message)
}
func (s *MessageService) SendSMS(
	req models.SMSRequest,
) error {

	content := ""

	switch req.Template {

	case "otp":

		otp := fmt.Sprintf(
			"%v",
			req.Data["otp"],
		)

		content = templates.OTPSMSTemplate(
			otp,
		)

	default:
		content = "No template found"
	}

	message := models.Message{
		To:       req.To,
		Content:  content,
		Template: req.Template,
	}

	return s.smsSender.Send(message)
}

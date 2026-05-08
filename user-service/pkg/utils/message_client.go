package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func sendEmailTemplate(
	to string,
	subject string,
	template string,
	data map[string]string,
) error {
	payload := map[string]interface{}{
		"to":       to,
		"subject":  subject,
		"template": template,
		"data":     data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal email payload: %w", err)
	}

	resp, err := http.Post(
		"http://message-service:3002/email/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("message-service email failed: status=%d", resp.StatusCode)
	}

	return nil
}

func SendOTPEmail(
	name string,
	email string,
	otp string,
) error {

	return sendEmailTemplate(
		email,
		"Your OTP Code",
		"otp",
		map[string]string{
			"name": name,
			"otp":  otp,
		},
	)
<<<<<<< Updated upstream

	return err
=======
}

func SendVerificationEmail(
	name string,
	email string,
	link string,
) error {

	return sendEmailTemplate(
		email,
		"Verify Your Email",
		"verify_email",
		map[string]string{
			"name": name,
			"link": link,
		},
	)
}

func SendForgotPasswordVerificationEmail(
	name string,
	email string,
	verifyLink string,
) error {
	return sendEmailTemplate(
		email,
		"Verify Password Reset Request",
		"forgot_password_verify",
		map[string]string{
			"name": name,
			"link": verifyLink,
		},
	)
}

func SendResetPasswordEmail(
	name string,
	email string,
	resetLink string,
) error {
	return sendEmailTemplate(
		email,
		"Reset Your Password",
		"reset_password",
		map[string]string{
			"name": name,
			"link": resetLink,
		},
	)
}

func SendUpdatePasswordEmail(
	name string,
	email string,
	updateLink string,
) error {
	return sendEmailTemplate(
		email,
		"Update Password",
		"update_password",
		map[string]string{
			"name": name,
			"link": updateLink,
		},
	)
>>>>>>> Stashed changes
}
package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func SendWelcomeEmail(
	name string,
	email string,
) error {

	return sendEmail(
		email,
		"Welcome to Our Platform",
		"welcome",
		map[string]string{
			"name": name,
		},
	)
}

func SendOTPEmail(name, email, otp string) error {
	return sendEmail(
		email,
		"Your OTP Code",
		"otp",
		map[string]string{
			"name": name,
			"otp":  otp,
		},
	)
}

func SendVerificationEmail(name, email, link string) error {
	return sendEmail(
		email,
		"Verify Your Email",
		"verify_email",
		map[string]string{
			"name": name,
			"link": link,
		},
	)
}

func sendEmail(to, subject, template string, data map[string]string) error {
	payload := map[string]interface{}{
		"to":       to,
		"subject":  subject,
		"template": template,
		"data":     data,
	}

	jsonData, _ := json.Marshal(payload)

	resp, err := httpClient.Post(
		"http://message-service:3002/email/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("message-service returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

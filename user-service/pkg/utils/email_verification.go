package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendEmailVerification(email, name, token string) error {

	verifyLink := "http://user-service:3001/api/v1/user/verify-email?token=" + token

	payload := map[string]interface{}{
		"to":       email,
		"subject":  "Verify your email",
		"template": "email_verification",
		"data": map[string]string{
			"name": name,
			"link": verifyLink,
		},
	}

	jsonData, _ := json.Marshal(payload)

	_, err := http.Post(
		"http://message-service:3002/email/verify",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	return err
}
package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendOTPEmail(
	name string,
	email string,
	otp string,
) error {

	payload := map[string]interface{}{
		"to": email,
		"subject": "Your OTP Code",
		"template": "otp",
		"data": map[string]string{
			"name": name,
			"otp":  otp,
		},
	}

	jsonData, _ := json.Marshal(payload)

	_, err := http.Post(
		"http://message-service:3002/email/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	return err
}
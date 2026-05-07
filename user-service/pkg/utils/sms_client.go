package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendOTPSMS(
	phone string,
	otp string,
) error {

	payload := map[string]interface{}{
		"to": phone,
		"template": "otp",
		"data": map[string]string{
			"otp": otp,
		},
	}

	jsonData, _ := json.Marshal(payload)

	_, err := http.Post(
		"http://message-service:3002/sms/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	return err
}
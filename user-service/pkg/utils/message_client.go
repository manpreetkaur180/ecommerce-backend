package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendOTPEmail(
	name string,
	email string,
	otp string,
) error {

	payload := map[string]interface{}{
		"to":       email,
		"subject":  "Your OTP Code",
		"template": "otp",
		"data": map[string]string{
			"name": name,
			"otp":  otp,
		},
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		return err
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("message service failed")
	}

	return nil
}
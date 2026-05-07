package clients

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendWelcomeEmail(
	name string,
	email string,
) error {

	payload := map[string]interface{}{
		"to": email,
		"subject": "Welcome to Our Platform",
		"template": "welcome",
		"data": map[string]string{
			"name": name,
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
package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func NotifyUserRegistered(
	name string,
	email string,
	phone string,
	token string,
) error {

	payload := map[string]string{
	"name":  name,
	"email": email,
	"phone": phone,
	"token": token,
}

	jsonData, _ := json.Marshal(payload)

	_, err := http.Post(
		"http://notification-service:3003/notify/user-registered",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	return err
}
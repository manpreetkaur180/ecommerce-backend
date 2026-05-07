package models

type SMSRequest struct {
	To       string                 `json:"to"`
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}
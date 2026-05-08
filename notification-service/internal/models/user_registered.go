package models

type UserRegisteredEvent struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type OTPEmailEvent struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type VerificationEmailEvent struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Link  string `json:"link"`
}

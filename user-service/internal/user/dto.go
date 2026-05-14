package user

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type OTPLoginRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	OTP   string `json:"otp"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type UpdatePasswordRequest struct {
	Email string `json:"email"`
}

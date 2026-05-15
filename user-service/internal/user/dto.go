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



type AddAddressRequest struct {
	FullName     string `json:"full_name" validate:"required,min=2,max=100"`
	Phone        string `json:"phone" validate:"required,phone"`
	AddressLine1 string `json:"address_line1" validate:"required,min=5,max=200"`
	AddressLine2 string `json:"address_line2,omitempty" validate:"max=200"`
	Landmark     string `json:"landmark,omitempty" validate:"max=100"`
	City         string `json:"city" validate:"required,min=2,max=100"`
	State        string `json:"state" validate:"required,min=2,max=100"`
	Country      string `json:"country" validate:"required,min=2,max=100"`
	Pincode      string `json:"pincode" validate:"required,len=6,numeric"`
	AddressType  string `json:"address_type" validate:"omitempty,oneof=home work other"`
	IsDefault    bool   `json:"is_default"`
}
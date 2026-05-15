package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"github.com/go-playground/validator/v10"
)
var validate = validator.New()


func ValidateStruct(i any) error {
	err := validate.Struct(i)
	if err == nil {
		return nil
	}

	var msgs []string

	for _, e := range err.(validator.ValidationErrors) {
		msgs = append(msgs, fieldError(e))
	}

	return errors.New(strings.Join(msgs, ", "))
}
func fieldError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be valid email"
	case "min":
		return e.Field() + " is too short"
	case "max":
		return e.Field() + " is too long"
	case "len":
		return e.Field() + " must be " + e.Param() + " characters"
	case "numeric":
		return e.Field() + " must be numeric"
	case "eqfield":
		return e.Field() + " must match " + e.Param()
	case "oneof":
		return e.Field() + " must be one of: " + e.Param()
	default:
		return e.Field() + " is invalid"
	}
}
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func NormalizePhone(phone string) string {

	phone = strings.TrimSpace(phone)

	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	return phone
}

// -------- EMAIL --------
func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(regex, email)

	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

// -------- PHONE --------
func ValidatePhone(phone string) error {
	if phone == "" {
		return nil // optional
	}

	regex := `^\+[1-9]\d{7,14}$`
	matched, _ := regexp.MatchString(regex, phone)

	if !matched {
		return errors.New("invalid phone number")
	}
	return nil
}

// -------- PASSWORD --------
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var hasUpper, hasLower, hasNumber bool

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New("password must contain uppercase, lowercase and number")
	}

	return nil
}
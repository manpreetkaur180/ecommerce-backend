package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

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
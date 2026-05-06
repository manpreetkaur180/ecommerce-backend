package utils

import (
	"errors"
	"regexp"
	"unicode"
)

// Email validation
func ValidateEmail(email string) error {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(regex, email)

	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

// Phone validation (simple 10-digit)
func ValidatePhone(phone string) error {
	if phone == "" {
		return nil // optional
	}

	regex := `^[0-9]{10}$`
	matched, _ := regexp.MatchString(regex, phone)

	if !matched {
		return errors.New("invalid phone number")
	}
	return nil
}

// Password validation
func ValidatePassword(password string) error {
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
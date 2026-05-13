package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

type PasswordResetClaims struct {
	UserID  uint   `json:"user_id"`
	Purpose string `json:"purpose"`

	jwt.RegisteredClaims
}

const passwordResetPurpose = "password_reset"

func ValidateJWTSecret() error {
	if os.Getenv("JWT_SECRET") == "" {
		return errors.New("JWT_SECRET is required")
	}

	return nil
}

// existing random token generator
func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// generate jwt
func GenerateJWT(
	userID uint,
	role string,
) (string, error) {

	expHours := 24

	if env := os.Getenv("JWT_EXPIRES_HOURS"); env != "" {
		if parsed, err := strconv.Atoi(env); err == nil {
			expHours = parsed
		}
	}

	claims := JWTClaims{
		UserID: userID,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Duration(expHours) * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
}

// parse jwt
func ParseJWT(tokenString string) (*JWTClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)

	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func GeneratePasswordResetJWT(userID uint) (string, error) {
	claims := PasswordResetClaims{
		UserID:  userID,
		Purpose: passwordResetPurpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)
}

func ParsePasswordResetJWT(tokenString string) (*PasswordResetClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&PasswordResetClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*PasswordResetClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid reset token")
	}

	if claims.Purpose != passwordResetPurpose {
		return nil, errors.New("invalid reset token purpose")
	}

	return claims, nil
}

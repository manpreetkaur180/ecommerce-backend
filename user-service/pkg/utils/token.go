package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	IsSeller bool   `json:"is_seller"`

	jwt.RegisteredClaims
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
	isSeller bool,
) (string, error) {

	expHours := 24

	if env := os.Getenv("JWT_EXPIRES_HOURS"); env != "" {
		if parsed, err := strconv.Atoi(env); err == nil {
			expHours = parsed
		}
	}

	claims := JWTClaims{
		UserID:   userID,
		Role:     role,
		IsSeller: isSeller,

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
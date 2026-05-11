package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)
type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IsSeller  bool   `json:"is_seller"`

	jwt.RegisteredClaims
}

func RequireAuth() fiber.Handler {

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Bearer TOKEN
		tokenString := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		if tokenString == authHeader {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.ParseWithClaims(
			tokenString,
			&JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			},
		)

		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		claims, ok := token.Claims.(*JWTClaims)

		if !ok || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		// store in context
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		c.Locals("is_seller", claims.IsSeller)

		return c.Next()
	}
}	
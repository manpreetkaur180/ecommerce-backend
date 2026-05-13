package middleware

import (
	"os"
	"strings"

	"cart-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponse(c, 401, "missing authorization header")
		}

		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			return utils.ErrorResponse(c, 401, "invalid authorization format")
		}

		tokenString := splitToken[1]

		token, err := jwt.ParseWithClaims(
			tokenString,
			&JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		if err != nil || !token.Valid {
			return utils.ErrorResponse(c, 401, "invalid or expired token")
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || claims.UserID == 0 || claims.Role == "" {
			return utils.ErrorResponse(c, 401, "invalid token claims")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RequireRoles(allowed ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		role, ok := c.Locals("role").(string)
		if !ok || role == "" {
			return utils.ErrorResponse(c, 403, "role not found")
		}

		for _, r := range allowed {
			if r == role {
				return c.Next()
			}
		}

		return utils.ErrorResponse(c, 403, "access denied")
	}
}

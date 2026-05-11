	package middleware

import (
	"strings"
	"user-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireAuth() fiber.Handler {

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return utils.ErrorResponse(
				c,
				401,
				"missing authorization header",
			)
		}

		splitToken := strings.Split(authHeader, " ")

		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			return utils.ErrorResponse(
				c,
				401,
				"invalid authorization format",
			)
		}

		tokenString := splitToken[1]

		claims, err := utils.ParseJWT(tokenString)

		if err != nil {
			return utils.ErrorResponse(
				c,
				401,
				"invalid or expired token",
			)
		}

		// store user data in request context
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("is_seller", claims.IsSeller)

		return c.Next()
	}
}
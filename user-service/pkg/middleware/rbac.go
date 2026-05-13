package middleware

import (
	"user-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)

		if !ok || role == "" {
			return utils.ErrorResponse(
				c,
				403,
				"role not found",
			)
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return utils.ErrorResponse(
			c,
			403,
			"access denied",
		)
	}
}

package middleware

import (
	"user-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

const adminRole = "admin"

func RequireAdmin() fiber.Handler {

	return func(c *fiber.Ctx) error {

		role, ok := c.Locals("role").(string)

		if !ok || role != adminRole {
			return utils.ErrorResponse(
				c,
				403,
				"admin access required",
			)
		}

		return c.Next()
	}
}

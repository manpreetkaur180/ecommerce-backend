package middleware

import (
	"user-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireSeller() fiber.Handler {

	return func(c *fiber.Ctx) error {

		role, ok := c.Locals("role").(string)

		if !ok || role != "seller" {
			return utils.ErrorResponse(
				c,
				403,
				"seller access required",
			)
		}

		return c.Next()
	}
}

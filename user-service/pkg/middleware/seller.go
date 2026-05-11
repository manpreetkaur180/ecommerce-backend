package middleware

import (
	"user-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireSeller() fiber.Handler {

	return func(c *fiber.Ctx) error {

		isSeller, ok := c.Locals("is_seller").(bool)

		if !ok || !isSeller {
			return utils.ErrorResponse(
				c,
				403,
				"seller access required",
			)
		}

		return c.Next()
	}
}
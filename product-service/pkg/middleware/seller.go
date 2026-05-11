package middleware

import "github.com/gofiber/fiber/v2"

func RequireSeller() fiber.Handler {

	return func(c *fiber.Ctx) error {

		isSeller, ok := c.Locals("is_seller").(bool)

		if !ok || !isSeller {

			return c.Status(403).JSON(fiber.Map{
				"error": "seller access required",
			})
		}

		return c.Next()
	}
}
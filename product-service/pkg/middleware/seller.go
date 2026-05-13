package middleware

import "github.com/gofiber/fiber/v2"

func RequireSeller() fiber.Handler {

	return func(c *fiber.Ctx) error {

		role, ok := c.Locals("role").(string)

		if !ok || role != "seller" {

			return c.Status(403).JSON(fiber.Map{
				"error": "seller access required",
			})
		}

		return c.Next()
	}
}

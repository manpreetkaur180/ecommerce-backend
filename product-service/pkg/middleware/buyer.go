package middleware

import "github.com/gofiber/fiber/v2"

func RequireBuyer() fiber.Handler {

	return func(c *fiber.Ctx) error {

		role, ok := c.Locals("role").(string)

		if !ok || role != "buyer" {

			return c.Status(403).JSON(fiber.Map{
				"error": "buyer access required",
			})
		}

		return c.Next()
	}
}

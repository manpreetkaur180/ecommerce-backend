package middleware

import "github.com/gofiber/fiber/v2"

func RequireBuyer() fiber.Handler {

	return func(c *fiber.Ctx) error {

		isSeller, ok := c.Locals("is_seller").(bool)

		// sellers are NOT allowed
		if !ok || isSeller {

			return c.Status(403).JSON(fiber.Map{
				"error": "buyer access required",
			})
		}

		return c.Next()
	}
}
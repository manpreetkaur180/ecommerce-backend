package middleware

import "github.com/gofiber/fiber/v2"

func RequireRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)

		if !ok || role == "" {
			return c.Status(403).JSON(fiber.Map{
				"error": "role not found",
			})
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "access denied",
		})
	}
}

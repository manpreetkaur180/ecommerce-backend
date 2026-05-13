package utils

import "github.com/gofiber/fiber/v2"

func SuccessResponse(
	c *fiber.Ctx,
	code int,
	message string,
	data interface{},
) error {
	return c.Status(code).JSON(fiber.Map{
		"message": message,
		"code":    code,
		"data":    data,
	})
}

func ErrorResponse(
	c *fiber.Ctx,
	code int,
	message string,
) error {
	return c.Status(code).JSON(fiber.Map{
		"error": message,
		"code":  code,
	})
}

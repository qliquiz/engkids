package errors

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func Handle(c *fiber.Ctx, err error) error {
	var ferr *fiber.Error
	if errors.As(err, &ferr) {
		return c.Status(ferr.Code).JSON(fiber.Map{"error": ferr.Message})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Внутренняя ошибка"})
}

package utils

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// ValidateStruct валидирует структуру на основе тегов validate
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func ParseAndValidate(c *fiber.Ctx, dst interface{}) error {
	if err := c.BodyParser(dst); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Невозможно обработать данные")
	}
	if err := ValidateStruct(dst); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return nil
}

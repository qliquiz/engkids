package utils

import (
	"github.com/go-playground/validator"
)

var validate = validator.New()

// ValidateStruct валидирует структуру на основе тегов validate
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
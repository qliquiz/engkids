package utils

import "github.com/go-playground/validator"

var Validate = validator.New()

func ValidateStruct(s interface{}) map[string]string {
	errs := map[string]string{}
	if err := Validate.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errs[err.Field()] = err.Tag()
		}
	}
	return errs
}

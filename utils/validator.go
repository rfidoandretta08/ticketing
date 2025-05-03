package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	// Validate the struct
	err := validate.Struct(s)
	if err != nil {
		return err
	}
	return nil
}

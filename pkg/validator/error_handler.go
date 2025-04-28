package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/yusuffugurlu/go-project/pkg/errors"
)

func ProcessValidationErrors(err error) *errors.AppError {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.NewBadRequest(err, "Invalid request format")
	}

	errorDetails := make(map[string]string)
	for _, e := range validationErrors {
		field := e.Field()
		errorDetails[field] = getValidationErrorMessage(e)
	}

	return errors.NewValidationError(err, errorDetails)
}

func getValidationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is less than minimum"
	case "max":
		return "Value is greater than maximum"
	default:
		return "Invalid value"
	}
}

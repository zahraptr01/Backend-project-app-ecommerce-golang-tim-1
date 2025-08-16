package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// validator inline message
func ValidateInput(data interface{}) (string, error) {
	// Create new validation and check the struct
	validate = validator.New()
	err := validate.Struct(data)

	if err != nil {
		// Handle invalid validation errors
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return "", nil
		}

		// Collect validation errors
		var errors []string
		for _, e := range err.(validator.ValidationErrors) {
			var message string
			if e.Tag() == "email" {
				message = "Please input correct email format"
			} else {
				message = fmt.Sprintf("%s must %s", e.Field(), e.Tag())
			}
			errors = append(errors, message)
		}
		return fmt.Sprint(errors), err
	}

	return "", nil
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// validator object struct message
func ValidateData(data interface{}) ([]FieldError, error) {
	validate := validator.New()

	err := validate.Struct(data)
	if err == nil {
		return nil, nil
	}

	var errors []FieldError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			var message string
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				message = "Please enter a valid email format"
			case "gte":
				message = fmt.Sprintf("%s must be a non-negative number", err.Field())
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
			case "eqfield":
				message = fmt.Sprintf("%s must match %s", err.Field(), err.Param())
			default:
				message = fmt.Sprintf("%s is invalid", err.Field())
			}

			errors = append(errors, FieldError{
				Field:   err.Field(),
				Message: message,
			})
		}
		return errors, err
	}

	// Fallback: return original error if not a validation error
	return nil, err
}

package service

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"avito-test-task/internal/lib/api/msg"
)

// ValidationError is a custom error type for validation errors.
type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

const (
	validateRequired = "required"
)

// ValidationErr returns a custom error message for validation errors.
func ValidationErr(errs validator.ValidationErrors) ValidationError {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case validateRequired:
			errMsgs = append(errMsgs, msg.ErrRequiredField(err.Field()))
		default:
			errMsgs = append(errMsgs, msg.ErrInvalidField(err.Field()))
		}
	}

	return ValidationError(strings.Join(errMsgs, ", "))
}

package msg

import "fmt"

func ErrRequiredField(field string) string {
	return fmt.Sprintf("field %s is a required field", field)
}

func ErrInvalidField(field string) string {
	return fmt.Sprintf("field %s is not valid", field)
}

func ErrInvalidFieldType(field, got, expected string) string {
	return fmt.Sprintf("expected type %s for field %s but got %s", expected, field, got)
}

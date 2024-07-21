package msg

import "fmt"

// ErrRequiredField returns a formatted string, indicating that field is a required field.
func ErrRequiredField(field string) string {
	return fmt.Sprintf("field %s is a required field", field)
}

// ErrInvalidField returns a formatted string, indicating that field is not valid.
func ErrInvalidField(field string) string {
	return fmt.Sprintf("field %s is not valid", field)
}

// ErrInvalidFieldType returns a formatted string, indicating that field has invalid field type.
func ErrInvalidFieldType(field, got, expected string) string {
	return fmt.Sprintf("expected type %s for field %s but got %s", expected, field, got)
}

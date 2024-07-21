package er

import (
	"strings"
)

const defaultSeparator = ", "

// Unwrap unwraps error err and returns a string, formed from all the errors from err, separated by comma.
func Unwrap(err error) string {
	return UnwrapSep(err, defaultSeparator)
}

// UnwrapSep unwraps error err and returns a string, formed from all the errors from err, separated by sep.
func UnwrapSep(err error, sep string) string {
	u, ok := err.(interface {
		Unwrap() []error
	})
	if !ok {
		return err.Error()
	}
	unwrapped := u.Unwrap()
	errs := make([]string, len(unwrapped))
	for i, e := range unwrapped {
		errs[i] = e.Error()
	}

	return strings.Join(errs, sep)
}

package er

import (
	"strings"
)

const defaultSeparator = ", "

func Unwrap(err error) string {
	return UnwrapSep(err, defaultSeparator)
}

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

package msg

import "fmt"

const (
	APIUnknownErr     = "unknown error"
	APIInternalErr    = "internal error"
	APIInvalidRequest = "invalid request"
	APIEmptyRequest   = "empty request"
	APINotAuthorized  = "only authorized users can access this resource"
	APIForbidden      = "forbidden"
)

func APIEmptyParameter(pName string) string {
	return fmt.Sprintf("empty parameter: %s", pName)
}

func APIUnacceptableFormat(pName string) string {
	return fmt.Sprintf("unacceptable format for parameter %s", pName)
}

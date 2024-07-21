package api

// Response is the general response struct that is being returned by http handlers.
type Response struct {
	Error string `json:"error,omitempty"`
}

// OkResponse is an empty response with no Error.
func OkResponse() Response {
	return Response{}
}

// ErrResponse is the response with some error err.
func ErrResponse(err string) Response {
	return Response{Error: err}
}

package api

type Response struct {
	Error string `json:"error,omitempty"`
}

func OkResponse() Response {
	return Response{}
}

func ErrResponse(err string) Response {
	return Response{Error: err}
}

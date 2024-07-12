package jsn

import (
	"avito-test-task/internal/lib/api/msg"
	"avito-test-task/internal/lib/logger/sl"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

// DecodingError is an error type that is returned when an error occurs during decoding.
type DecodingError string

func (e DecodingError) Error() string {
	return string(e)
}

// EncodeResponse writes the response to the http.ResponseWriter.
// If an error occurs, it logs the error and writes the default
// error message (msg.APIUnknownErr) to the http.ResponseWriter.
func EncodeResponse(w http.ResponseWriter, statusCode int, response any, log *slog.Logger) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil && !errors.Is(err, http.ErrBodyNotAllowed) {
		log.Error("failed to encode response", sl.Err(err))
		http.Error(w, msg.APIUnknownErr, http.StatusInternalServerError)
	}
}

// DecodeRequest reads the request body and decodes it into the provided request object.
// If an error occurs, it logs the error and returns a DecodingError.
func DecodeRequest[T any](r *http.Request, req *T, log *slog.Logger) error {
	if err := json.NewDecoder(r.Body).Decode(req); errors.Is(err, io.EOF) {
		log.Error("empty request body", sl.Err(err))
		return DecodingError(msg.APIEmptyRequest)
	} else if unmarshalErr := new(json.UnmarshalTypeError); errors.As(err, &unmarshalErr) {
		log.Error("failed to decode request", sl.Err(err))
		return DecodingError(msg.ErrInvalidFieldType(unmarshalErr.Field, unmarshalErr.Value, unmarshalErr.Type.String()))
	} else if err != nil {
		log.Error("failed to decode request", sl.Err(err))
		return DecodingError(msg.APIInvalidRequest)
	}

	return nil
}

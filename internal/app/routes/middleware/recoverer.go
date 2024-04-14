package middleware

import (
	"avito-test-task/internal/lib/api"
	"avito-test-task/internal/lib/api/jsn"
	"avito-test-task/internal/lib/api/msg"
	"avito-test-task/internal/lib/logger/sl"
	"fmt"
	"log/slog"
	"net/http"
)

// NewRecovererMiddleware is a middleware that recovers from panics and returns
// a 500 Internal Server Error in such cases.
// It sets an error message in the response body.
func NewRecovererMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic occurred. recovered.", sl.Err(fmt.Errorf("%v", err)))
					w.Header().Set("Content-Type", "application/json")
					jsn.EncodeResponse(w, http.StatusInternalServerError, api.ErrResponse(msg.APIInternalErr), logger)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

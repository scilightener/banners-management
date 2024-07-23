package middleware

import (
	"net/http"

	"github.com/google/uuid"

	"banners-management/internal/lib/api"
)

// RequestIDMiddleware is a middleware that adds a unique request ID to each request.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		r = api.SetRequestID(r, requestID)
		next.ServeHTTP(w, r)
	})
}

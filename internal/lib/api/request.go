package api

import (
	"context"
	"net/http"
)

const (
	RequestIDKey = "request-id"
	RoleKey      = "role"
)

// RequestID returns request id, associated with the given request.
func RequestID(r *http.Request) string {
	return ctxValue(r.Context(), RequestIDKey)
}

// SetRequestID return a request with the given request id.
// Request id can be retrieved with RequestID function.
func SetRequestID(r *http.Request, requestID string) *http.Request {
	ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
	return r.WithContext(ctx)
}

// UserRole returns user role, associated with the user, making request.
func UserRole(r *http.Request) string {
	return ctxValue(r.Context(), RoleKey)
}

// SetUserRole return a context with the given user role.
// User role can be retrieved with UserRole function.
func SetUserRole(r *http.Request, role string) *http.Request {
	ctx := context.WithValue(r.Context(), RoleKey, role)
	return r.WithContext(ctx)
}

// ctxValue returns a value from the context by the given key.
func ctxValue(ctx context.Context, key string) string {
	if value := ctx.Value(key); value != nil {
		return value.(string)
	}

	return ""
}

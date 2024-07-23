package middleware

import (
	"banners-management/internal/lib/api"
	"banners-management/internal/lib/api/jsn"
	"banners-management/internal/lib/api/msg"
	"banners-management/internal/lib/jwt"
	"banners-management/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"strings"
)

const Authorization = "Authorization"

// NewAuthorizationMiddleware creates a new authorization middleware.
// It checks the Authorization header for a valid JWT token.
// If the token is valid, it extracts the role from it and adds it to the request context.
func NewAuthorizationMiddleware(logger *slog.Logger, manager *jwt.Manager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(Authorization)
			if token == "" {
				logger.Info("nothing in Authorization header")
				jsn.EncodeResponse(w, http.StatusUnauthorized, api.ErrResponse(msg.APINotAuthorized), logger)
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			err := manager.VerifyToken(token)
			if err != nil {
				logger.Info("invalid jwt token", sl.Err(err))
				jsn.EncodeResponse(w, http.StatusUnauthorized, api.ErrResponse(msg.APINotAuthorized), logger)
				return
			}

			role, err := manager.GetRole(token)
			if err != nil {
				logger.Info("failed to get role from token", sl.Err(err))
				jsn.EncodeResponse(w, http.StatusUnauthorized, api.ErrResponse(msg.APINotAuthorized), logger)
				return
			}

			r = api.SetUserRole(r, role)

			next.ServeHTTP(w, r)
		})
	}
}

// EnsureAdmin returns new http.Handler that checks if the incoming request authorized with admin role,
// and if so, gives access to the calling endpoint, otherwise returns 403 Forbidden status code response.
func EnsureAdmin(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := api.UserRole(r)
		if role != "admin" {
			jsn.EncodeResponse(w, http.StatusForbidden, api.ErrResponse(msg.APIForbidden), logger)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package middleware

import (
	"banners-management/internal/lib/api"
	"log/slog"
	"net/http"
	"time"
)

// NewLoggingMiddleware creates a new logging middleware.
// It logs the request method, path, remote address, user agent, and request ID, response status code and its duration.
func NewLoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String(api.RequestIDKey, api.RequestID(r)),
			)

			log.Info("request started")

			wrw := &wrappedResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			t1 := time.Now()

			next.ServeHTTP(wrw, r)

			log.Info("request completed",
				slog.Int("status", wrw.statusCode),
				slog.String("duration", time.Since(t1).String()),
			)
		})
	}
}

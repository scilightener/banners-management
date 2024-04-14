package middleware

import "net/http"

// ContentTypeJSONMiddleware sets the Content-Type header to application/json for every request.
func ContentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

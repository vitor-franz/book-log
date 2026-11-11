package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRootMiddleware returns middleware that sets the common JSON response headers.
func NewRootMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("X-Content-Type-Options", "nosniff")

			next.ServeHTTP(w, r)
		})
	}
}

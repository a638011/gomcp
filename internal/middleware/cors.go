package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// CORSMiddleware creates a CORS middleware with specified configuration
func CORSMiddleware(enabled bool, origins []string, credentials bool, methods []string, headers []string) func(http.Handler) http.Handler {
	if !enabled {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   methods,
		AllowedHeaders:   headers,
		AllowCredentials: credentials,
		MaxAge:           300,
	})
}

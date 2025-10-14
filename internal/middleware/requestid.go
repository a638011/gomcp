package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey contextKey = "request_id"
	// RequestIDHeader is the HTTP header name for request ID
	RequestIDHeader = "X-Request-ID"
)

// RequestID middleware adds a unique request ID to each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists in header
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			// Generate new request ID
			requestID = uuid.New().String()
		}

		// Add request ID to context
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		// Add request ID to response header
		w.Header().Set(RequestIDHeader, requestID)

		next.ServeHTTP(w, r)
	})
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

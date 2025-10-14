package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// SessionMiddleware creates a session middleware
func SessionMiddleware(sessionSecret string) func(http.Handler) http.Handler {
	store := sessions.NewCookieStore([]byte(sessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false, // Set to true in production with HTTPS
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Session is available via gorilla/sessions.Get(r, "mcp_session")
			next.ServeHTTP(w, r)
		})
	}
}

package middleware

import (
	"net/http"

	"github.com/souravendra/distributed-rate-limiter-go-crispy-tribble/limiter"
)

// --- Middleware Pattern ---
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := "test-client" // r.RemoteAddr
		if !limiter.GetLimiter().Strategy.Allow(key) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

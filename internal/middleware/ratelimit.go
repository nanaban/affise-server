package middleware

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// RateLimit represents rate limit middleware.
type RateLimit struct {
	limiter *rate.Limiter
}

// NewRateLimit creates new instance of rate limit middleware.
func NewRateLimit(interval time.Duration, limit int) *RateLimit {
	return &RateLimit{
		limiter: rate.NewLimiter(rate.Every(interval), limit),
	}
}

// Handle handles rate limit middleware.
func (rl *RateLimit) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.limiter.Allow() {
			status := http.StatusTooManyRequests
			http.Error(w, http.StatusText(status), status)
			return
		}

		next.ServeHTTP(w, r)
	})
}

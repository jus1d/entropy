package ratelimit

import (
	"sync"

	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
)

// TokenBucket returns a rate limiting middleware using the token bucket algorithm.
// r is the refill rate (tokens per second), burst is the bucket capacity.
// Each IP gets its own bucket that starts full and refills at rate r.
func TokenBucket(r rate.Limit, burst int) echo.MiddlewareFunc {
	var (
		mu       sync.Mutex
		visitors = make(map[string]*rate.Limiter)
	)

	return middleware(func(ip string) bool {
		mu.Lock()
		l, ok := visitors[ip]
		if !ok {
			l = rate.NewLimiter(r, burst)
			visitors[ip] = l
		}
		mu.Unlock()

		return l.Allow()
	})
}

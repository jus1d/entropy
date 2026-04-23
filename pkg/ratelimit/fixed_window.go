package ratelimit

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type fixedWindowEntry struct {
	count int
	start time.Time
}

// FixedWindow returns a rate limiting middleware using the fixed window algorithm.
// limit is the max number of requests allowed per window duration.
// The counter resets at each window boundary.
func FixedWindow(limit int, window time.Duration) echo.MiddlewareFunc {
	var (
		mu       sync.Mutex
		visitors = make(map[string]*fixedWindowEntry)
	)

	return middleware(func(ip string) bool {
		now := time.Now()

		mu.Lock()
		defer mu.Unlock()

		e, ok := visitors[ip]
		if !ok || now.Sub(e.start) >= window {
			visitors[ip] = &fixedWindowEntry{count: 1, start: now}
			return true
		}

		e.count++
		return e.count <= limit
	})
}

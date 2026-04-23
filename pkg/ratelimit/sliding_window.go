package ratelimit

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type slidingWindowEntry struct {
	prevCount int
	currCount int
	currStart time.Time
}

// SlidingWindow returns a rate limiting middleware using the sliding window algorithm.
// limit is the max number of requests allowed per window duration.
// Unlike FixedWindow, this smooths the boundary between windows by weighting
// the previous window's count based on how far into the current window we are.
func SlidingWindow(limit int, window time.Duration) echo.MiddlewareFunc {
	var (
		mu       sync.Mutex
		visitors = make(map[string]*slidingWindowEntry)
	)

	return middleware(func(ip string) bool {
		now := time.Now()

		mu.Lock()
		defer mu.Unlock()

		e, ok := visitors[ip]
		if !ok {
			visitors[ip] = &slidingWindowEntry{currCount: 1, currStart: now}
			return true
		}

		elapsed := now.Sub(e.currStart)
		if elapsed >= window {
			e.prevCount = e.currCount
			e.currCount = 0
			e.currStart = now
			elapsed = 0
		}

		// Weight previous window by how much of it still overlaps
		weight := 1.0 - float64(elapsed)/float64(window)
		estimate := float64(e.prevCount)*weight + float64(e.currCount)

		e.currCount++
		return estimate < float64(limit)
	})
}

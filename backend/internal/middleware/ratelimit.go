package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Cleanup old entries every minute
	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, times := range rl.requests {
		var valid []time.Time
		for _, t := range times {
			if now.Sub(t) < rl.window {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(rl.requests, key)
		} else {
			rl.requests[key] = valid
		}
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times := rl.requests[key]

	// Filter old requests in-place to avoid allocation
	validCount := 0
	for _, t := range times {
		if now.Sub(t) < rl.window {
			times[validCount] = t
			validCount++
		}
	}
	times = times[:validCount]

	if validCount >= rl.limit {
		rl.requests[key] = times
		return false
	}

	times = append(times, now)
	rl.requests[key] = times
	return true
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limit int, window time.Duration) fiber.Handler {
	limiter := NewRateLimiter(limit, window)

	return func(c *fiber.Ctx) error {
		key := c.IP()

		if !limiter.Allow(key) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"retry_after": window.Seconds(),
			})
		}

		return c.Next()
	}
}

// StrictRateLimitMiddleware for sensitive endpoints (lower limits)
func StrictRateLimitMiddleware() fiber.Handler {
	return RateLimitMiddleware(10, time.Minute)
}

// APIRateLimitMiddleware for general API endpoints
func APIRateLimitMiddleware() fiber.Handler {
	return RateLimitMiddleware(100, time.Minute)
}

// PublicRateLimitMiddleware for public endpoints like SVG/JSON
func PublicRateLimitMiddleware() fiber.Handler {
	return RateLimitMiddleware(60, time.Minute)
}

// EnforceJSONMiddleware ensures that the client accepts JSON responses
func EnforceJSONMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only check for non-GET requests or specific API routes
		if c.Method() == fiber.MethodGet {
			return c.Next()
		}

		accept := c.Get("Accept")
		if accept != "" && accept != "*/*" && !contains(accept, "application/json") {
			return c.Status(fiber.StatusNotAcceptable).JSON(fiber.Map{
				"error": "Only application/json is supported",
			})
		}
		return c.Next()
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr)
}

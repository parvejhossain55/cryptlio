package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
)

func getLimiter(ip string) *rate.Limiter {
	limiter, exists := limiters[ip]
	if !exists {
		// Default: 100 requests per second
		limiter = rate.NewLimiter(rate.Every(time.Second), 100)
		limiters[ip] = limiter
	}
	return limiter
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"retry_after": time.Second.String(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetRateLimiter returns a rate limiter for a given key (e.g., user_id for authenticated endpoints)
func GetRateLimiter(rps int) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := rate.NewLimiter(rate.Every(time.Second), rps)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

package middleware

import (
	"sync"

	"github.com/afandimsr/go-gin-api/internal/delivery/http/response"
	"github.com/afandimsr/go-gin-api/internal/domain/apperror"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit returns a middleware that limits the number of requests per IP address.
func RateLimitPerIP(r rate.Limit, b int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	var mu sync.Mutex

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(r, b)
			limiters[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			response.Error(c, 429, apperror.RateLimitExceeded, "Too many requests", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

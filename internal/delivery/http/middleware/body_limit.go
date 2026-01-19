package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimit returns a middleware that limits the size of the request body.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			maxBytes,
		)
		c.Next()
	}
}

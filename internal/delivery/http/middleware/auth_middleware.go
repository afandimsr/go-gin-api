package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/username/go-gin-api/internal/domain/apperror"
	"github.com/username/go-gin-api/internal/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(apperror.Unauthorized("authorization header required", nil))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(apperror.Unauthorized("invalid authorization format", nil))
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(parts[1])
		if err != nil {
			c.Error(apperror.Unauthorized("invalid or expired token", err))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

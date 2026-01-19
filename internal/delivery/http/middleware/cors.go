package middleware

import (
	"strings"
	"time"

	"github.com/afandimsr/go-gin-api/internal/config"

	"github.com/gin-contrib/cors"
)

// Cors returns a CORS configuration based on the environment variables
func Cors(co *config.Config) cors.Config {
	origins := strings.Split(co.CorsAllowedOrigins, ",")

	return cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

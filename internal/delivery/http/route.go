package http

import (
	"github.com/gin-gonic/gin"
	"github.com/username/go-gin-api/internal/delivery/http/handler"
	"github.com/username/go-gin-api/internal/delivery/http/middleware"
)

func RegisterRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
) {
	api := r.Group("/api/v1")

	// auth routes
	api.POST("/login", userHandler.Login)

	// health check
	api.GET("/health", healthHandler)

	// user routes
	users := api.Group("/users")
	{
		users.GET("", userHandler.GetUsers)
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)

		protected := users.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.PUT("/:id", userHandler.UpdateUser)
			protected.DELETE("/:id", userHandler.DeleteUser)
		}
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

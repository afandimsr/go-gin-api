package http

import (
	"github.com/afandimsr/go-gin-api/internal/delivery/http/handler"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
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
		protected := users.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.PUT("/:id", userHandler.UpdateUser)
			protected.DELETE("/:id", userHandler.DeleteUser)
			protected.GET("", userHandler.GetUsers)
			protected.POST("", userHandler.CreateUser)
			protected.GET("/:id", userHandler.GetUser)
		}
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

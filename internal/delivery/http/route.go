package http

import (
	handler "github.com/afandimsr/go-gin-api/internal/delivery/http/handler/user"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/middleware"
	"github.com/afandimsr/go-gin-api/internal/domain/user"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	ks user.KeycloakService,
) {

	api := r.Group("/api/v1")

	// auth routes
	api.POST("/login", userHandler.Login)
	api.GET("/auth/login", userHandler.OIDCLogin)
	api.GET("/auth/callback", userHandler.OIDCCallback)
	api.GET("/logout", userHandler.Logout)

	// health check
	api.GET("/health", healthHandler)

	// user routes
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(ks), middleware.AdminOnly())
	{
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.GET("", userHandler.GetUsers)
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id/change-password", userHandler.ChangePassword)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

package bootstrap

import (
	"github.com/gin-gonic/gin"
	httpDelivery "github.com/username/go-gin-api/internal/delivery/http"
	"github.com/username/go-gin-api/internal/delivery/http/handler"
)

func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	httpDelivery.RegisterRoutes(r, userHandler)
}

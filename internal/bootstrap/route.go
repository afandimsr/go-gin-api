package bootstrap

import (
	httpDelivery "github.com/afandimsr/go-gin-api/internal/delivery/http"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler) {
	httpDelivery.RegisterRoutes(r, userHandler)
}

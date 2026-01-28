package bootstrap

import (
	httpDelivery "github.com/afandimsr/go-gin-api/internal/delivery/http"
	handler "github.com/afandimsr/go-gin-api/internal/delivery/http/handler/user"
	"github.com/afandimsr/go-gin-api/internal/domain/user"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	ks user.KeycloakService,
) {
	httpDelivery.RegisterRoutes(r, userHandler, ks)
}

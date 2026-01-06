package bootstrap

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/username/go-gin-api/docs"
	"github.com/username/go-gin-api/internal/config"
	"github.com/username/go-gin-api/internal/delivery/http/handler"
	"github.com/username/go-gin-api/internal/delivery/http/middleware"
	"github.com/username/go-gin-api/internal/infrastructure/external"
	userRepo "github.com/username/go-gin-api/internal/infrastructure/persistent/mysql/repository"
	"github.com/username/go-gin-api/internal/pkg/jwt"
	userUC "github.com/username/go-gin-api/internal/usecase/user"
)

func Run() {
	cfg := config.Load()
	jwt.SetSecret(cfg.JWTSecret)

	db, err := config.NewMySQL(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	authClient := external.NewAuthClient(cfg.ClientAuthURL)
	userRepository := userRepo.NewUserRepo(db)
	userUsecase := userUC.New(userRepository, authClient)
	userHandler := handler.New(userUsecase)

	r := gin.Default()
	r.Use(middleware.ErrorHandler())

	RegisterRoutes(r, userHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Running on port", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}

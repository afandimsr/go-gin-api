package bootstrap

import (
	"log"

	_ "github.com/afandimsr/go-gin-api/docs"
	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/afandimsr/go-gin-api/internal/database"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/handler"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/middleware"
	"github.com/afandimsr/go-gin-api/internal/infrastructure/external"
	userRepo "github.com/afandimsr/go-gin-api/internal/infrastructure/persistent/mysql/repository"
	"github.com/afandimsr/go-gin-api/internal/pkg/jwt"
	userUC "github.com/afandimsr/go-gin-api/internal/usecase/user"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run() {
	cfg := config.Load()
	jwt.SetSecret(cfg.JWTSecret)

	db, err := database.NewDatabase(cfg.DB)
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

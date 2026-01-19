package bootstrap

import (
	"log"

	_ "github.com/afandimsr/go-gin-api/docs"
	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/afandimsr/go-gin-api/internal/database"
	handler "github.com/afandimsr/go-gin-api/internal/delivery/http/handler/user"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/middleware"
	"github.com/afandimsr/go-gin-api/internal/infrastructure/external"
	userRepo "github.com/afandimsr/go-gin-api/internal/infrastructure/persistent/mysql/repository"
	userPostgresRepo "github.com/afandimsr/go-gin-api/internal/infrastructure/persistent/postgres/repository"
	"github.com/afandimsr/go-gin-api/internal/pkg/jwt"
	userUC "github.com/afandimsr/go-gin-api/internal/usecase/user"
	"github.com/gin-contrib/cors"
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
	var userHandler *handler.UserHandler
	switch cfg.DB.Driver {
	case "mysql":
		userRepository := userRepo.NewUserRepo(db)
		userUsecase := userUC.New(userRepository, authClient)
		userHandler = handler.New(userUsecase)
	case "postgres":
		userRepository := userPostgresRepo.NewUserRepo(db)
		userUsecase := userUC.New(userRepository, authClient)
		userHandler = handler.New(userUsecase)
	default:
		log.Fatal("Unsupported database driver: " + cfg.DB.Driver)
	}

	r := gin.Default()
	r.Use(
		cors.New(middleware.Cors(cfg)),
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.SecureHeaders(),
		middleware.BodyLimit(2<<20), // 2MB
		middleware.RateLimitPerIP(10, 20),
		middleware.ErrorHandler(),
	)

	RegisterRoutes(r, userHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Running on port", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}

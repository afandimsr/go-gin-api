package bootstrap

import (
	"context"
	"log"

	_ "github.com/afandimsr/go-gin-api/docs"
	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/afandimsr/go-gin-api/internal/database"
	handler "github.com/afandimsr/go-gin-api/internal/delivery/http/handler/user"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/middleware"
	"github.com/afandimsr/go-gin-api/internal/infrastructure/external"
	userRepo "github.com/afandimsr/go-gin-api/internal/infrastructure/persistent/mysql/repository"
	userPostgresRepo "github.com/afandimsr/go-gin-api/internal/infrastructure/persistent/postgres/repository"
	s3infra "github.com/afandimsr/go-gin-api/internal/infrastructure/storage/s3"
	"github.com/afandimsr/go-gin-api/internal/pkg/jwt"
	userUC "github.com/afandimsr/go-gin-api/internal/usecase/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run() {
	// load config
	cfg := config.Load()
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		log.Println("Production mode")
	}

	// set jwt secret
	jwt.SetSecret(cfg.JWTSecret)

	// initialize database
	db, err := database.NewDatabase(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	// initialize storage service
	ctx := context.Background()

	publicS3, err := s3infra.New(ctx, cfg.S3["public"])
	if err != nil {
		log.Fatal("failed init public s3:", err)
	}

	// privateS3, err := s3infra.New(ctx, cfg.S3["private"])
	// if err != nil {
	// 	log.Fatal("failed init private s3:", err)
	// }

	publicStorage := s3infra.NewUploader(
		publicS3,
		cfg.S3["public"].Bucket,
	)

	// privateStorage := storageUC.New(
	// 	privateS3,
	// 	cfg.S3["private"].Bucket,
	// )

	// inject to usecase layer
	_ = publicStorage
	// _ = privateStorage

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
		middleware.ErrorHandler(cfg),
	)

	RegisterRoutes(r, userHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Running on port", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}

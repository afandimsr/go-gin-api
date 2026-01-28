package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/afandimsr/go-gin-api/docs"
	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/afandimsr/go-gin-api/internal/database"
	handler "github.com/afandimsr/go-gin-api/internal/delivery/http/handler/user"
	"github.com/afandimsr/go-gin-api/internal/delivery/http/middleware"
	"github.com/afandimsr/go-gin-api/internal/infrastructure/apm"
	"github.com/afandimsr/go-gin-api/internal/infrastructure/external"
	userRepo "github.com/afandimsr/go-gin-api/internal/infrastructure/persistent/mysql/repository"
	userPostgresRepo "github.com/afandimsr/go-gin-api/internal/infrastructure/persistent/postgres/repository"
	s3infra "github.com/afandimsr/go-gin-api/internal/infrastructure/storage/s3"
	"github.com/afandimsr/go-gin-api/internal/pkg/jwt"
	"github.com/afandimsr/go-gin-api/internal/pkg/oidc"
	userUC "github.com/afandimsr/go-gin-api/internal/usecase/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	keycloakService := external.NewKeycloakService(cfg.Keycloak)

	// OIDC Provider (optional depending on config)
	var oidcProvider *oidc.OIDCProvider
	if cfg.Keycloak.URL != "" {
		redirectURL := fmt.Sprintf("http://localhost:%s/api/v1/auth/callback", cfg.AppPort)
		op, err := oidc.NewOIDCProvider(context.Background(), cfg.Keycloak, redirectURL)
		if err != nil {
			log.Printf("Warning: failed to initialize OIDC provider: %v", err)
		} else {
			oidcProvider = op
		}
	}

	var userHandler *handler.UserHandler
	switch cfg.DB.Driver {
	case "mysql":
		userRepository := userRepo.NewUserRepo(db)
		userUsecase := userUC.New(userRepository, authClient, keycloakService)
		userHandler = handler.New(userUsecase, oidcProvider)
	case "postgres":
		userRepository := userPostgresRepo.NewUserRepo(db)
		userUsecase := userUC.New(userRepository, authClient, keycloakService)
		userHandler = handler.New(userUsecase, oidcProvider)
	default:
		log.Fatal("Unsupported database driver: " + cfg.DB.Driver)
	}

	// initialize APM
	apm.Init(cfg)

	r := gin.Default()
	r.Use(
		cors.New(middleware.Cors(cfg)),
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.SecureHeaders(),
		middleware.BodyLimit(2<<20), // 2MB
		middleware.RateLimitPerIP(10, 20),
		apm.GinMiddleware(), // Gin Elastic APM
		middleware.ErrorHandler(cfg),
	)

	RegisterRoutes(r, userHandler, keycloakService)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Ensure roles exist
	seedRoles(db)

	log.Println("Running on port", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}

func seedRoles(db *sql.DB) {
	roles := []string{"USER", "ADMIN"}
	for _, role := range roles {
		var id string
		err := db.QueryRow("SELECT id FROM roles WHERE name = ?", role).Scan(&id)
		if err != nil {
			log.Printf("[Seed] Role %s missing, creating...", role)
			_, err = db.Exec("INSERT INTO roles (id, name) VALUES (?, ?)", uuid.New().String(), role)
			if err != nil {
				log.Printf("[Seed] Failed to create role %s: %v", role, err)
			}
		}
	}
}

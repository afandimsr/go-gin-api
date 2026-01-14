package main

import (
	"log"

	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/afandimsr/go-gin-api/internal/database"
	"github.com/afandimsr/go-gin-api/internal/seeder"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg.DB)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	switch cfg.DB.Driver {
	case "postgres":
		seeder.SeedRolesPostgres(db)
		seeder.SeedAdminUserPostgres(db)
	case "mysql":
		seeder.SeedRolesMysql(db)
		seeder.SeedAdminUserMysql(db)
	default:
		log.Fatalf("unsupported db driver: %s", cfg.DB.Driver)
	}
}

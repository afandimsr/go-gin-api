package main

import (
	"database/sql"
	"log"

	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/afandimsr/go-gin-api/internal/seeder"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()

	dbUser := cfg.DB.User
	dbPassword := cfg.DB.Password
	dbName := cfg.DB.Name
	dbHost := cfg.DB.Host
	dbPort := cfg.DB.Port
	dbSSLMode := cfg.DB.SSLMode

	dbURL := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + dbSSLMode
	log.Printf("Connecting to DB at %s", dbURL)

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// Ensure connection
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	// Seed roles
	seeder.SeedRoles(db)
	// Seed the default admin user
	seeder.SeedAdminUser(db)
}

package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/afandimsr/go-gin-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgres(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// ✅ CONNECTION POOLING
	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetMaxIdleConns(cfg.MaxIdle)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	stats := db.Stats()
	fmt.Println(stats.OpenConnections)
	fmt.Println(stats.InUse)
	fmt.Println(stats.Idle)

	return db, db.Ping()
}

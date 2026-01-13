package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/afandimsr/go-gin-api/internal/config"
)

func NewMySQL(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("mysql", dsn)
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

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

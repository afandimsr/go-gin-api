package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/afandimsr/go-gin-api/internal/config"
	_ "github.com/denisenkom/go-mssqldb"
)

func NewSQLServer(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("sqlserver", dsn)
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

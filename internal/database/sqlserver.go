package database

import (
	"database/sql"
	"fmt"

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

	db.SetMaxOpenConns(cfg.MaxOpen)
	db.SetMaxIdleConns(cfg.MaxIdle)

	return db, db.Ping()
}

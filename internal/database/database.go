package database

import (
	"database/sql"
	"fmt"

	"github.com/afandimsr/go-gin-api/internal/config"
)

func NewDatabase(cfg config.DBConfig) (*sql.DB, error) {
	switch cfg.Driver {
	case "postgres":
		return NewPostgres(cfg)
	case "mysql":
		return NewMySQL(cfg)
	case "sqlserver":
		return NewSQLServer(cfg)
	default:
		return nil, fmt.Errorf("unsupported db driver: %s", cfg.Driver)
	}
}

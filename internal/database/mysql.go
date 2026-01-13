package database

import (
	"database/sql"

	"github.com/afandimsr/go-gin-api/internal/config"
	"github.com/go-sql-driver/mysql"
)

func NewMySQL(cfg config.DBConfig) (*sql.DB, error) {
	dbCfg := mysql.Config{
		User:      cfg.User,
		Passwd:    cfg.Password,
		Net:       "tcp",
		Addr:      cfg.Host + ":" + cfg.Port,
		DBName:    cfg.Name,
		ParseTime: true,
	}

	return sql.Open("mysql", dbCfg.FormatDSN())
}

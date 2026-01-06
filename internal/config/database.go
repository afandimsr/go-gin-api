package config

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func NewMySQL(cfg DBConfig) (*sql.DB, error) {
	dbCfg := mysql.Config{
		User:      cfg.User,
		Passwd:    cfg.Pass,
		Net:       "tcp",
		Addr:      cfg.Host + ":" + cfg.Port,
		DBName:    cfg.Name,
		ParseTime: true,
	}

	return sql.Open("mysql", dbCfg.FormatDSN())
}

package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName       string
	AppPort       string
	AppEnv        string
	JWTSecret     string
	ClientAuthURL string

	DB DBConfig
}

type DBConfig struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	MaxOpen  int
	MaxIdle  int
}

func Load() *Config {
	// Load .env (ignore error in production)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config:", err)
	}

	cfg := &Config{
		AppName:       getEnv("APP_NAME", "go-app"),
		AppPort:       getEnv("APP_PORT", "8080"),
		AppEnv:        getEnv("APP_ENV", "development"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret"),
		ClientAuthURL: getEnv("CLIENT_AUTH_URL", ""),

		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
			Driver:   getEnv("DB_DRIVER", "mysql"),
			SSLMode:  getEnv("DB_SSL_MODE", ""),
			MaxOpen:  getEnvInt("DB_MAX_OPEN", 20),
			MaxIdle:  getEnvInt("DB_MAX_IDLE", 10),
		},
	}

	validate(cfg)
	return cfg
}

func getEnv(key string, defaultVal string) string {
	if val := viper.GetString(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := viper.GetInt(key); val != 0 {
		return val
	}
	return defaultVal
}

func validate(cfg *Config) {
	if cfg.DB.Name == "" {
		log.Fatal("DB_NAME is required")
	}
}

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
	Host string
	Port string
	User string
	Pass string
	Name string
}

func Load() *Config {
	// Load .env (ignore error in production)
	_ = godotenv.Load()

	cfg := &Config{
		AppName:       getEnv("APP_NAME", "go-app"),
		AppPort:       getEnv("APP_PORT", "8080"),
		AppEnv:        getEnv("APP_ENV", "production"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret"),
		ClientAuthURL: getEnv("CLIENT_AUTH_URL", ""),

		DB: DBConfig{
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "3306"),
			User: getEnv("DB_USER", "root"),
			Pass: getEnv("DB_PASS", ""),
			Name: getEnv("DB_NAME", ""),
		},
	}

	validate(cfg)
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func validate(cfg *Config) {
	if cfg.DB.Name == "" {
		log.Fatal("DB_NAME is required")
	}
}

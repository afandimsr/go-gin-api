package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName            string
	AppPort            string
	AppEnv             string
	JWTSecret          string
	ClientAuthURL      string
	CorsAllowedOrigins string

	DB         DBConfig
	Keycloak   KeycloakConfig
	S3         map[string]S3Config `mapstructure:"s3"`
	ElasticApm ElasticApmConfig
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

type KeycloakConfig struct {
	URL           string
	Realm         string
	ClientID      string
	ClientSecret  string
	AdminUser     string
	AdminPassword string
}

type ElasticApmConfig struct {
	ServerURL        string
	ServiceName      string
	Environment      string
	SecretToken      string
	VerifyServerCert bool
	ServiceVersion   string
}

func Load() *Config {
	// Load .env (ignore error in production)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	//
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Failed to read config:", err)
	}

	cfg := &Config{
		AppName:            getEnv("APP_NAME", "go-app"),
		AppPort:            getEnv("APP_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "development"),
		JWTSecret:          getEnv("JWT_SECRET", "default-secret"),
		ClientAuthURL:      getEnv("CLIENT_AUTH_URL", ""),
		CorsAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),

		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "postgres"),
			Driver:   getEnv("DB_DRIVER", "postgres"),
			SSLMode:  getEnv("DB_SSL_MODE", ""),
			MaxOpen:  getEnvInt("DB_MAX_OPEN", 20),
			MaxIdle:  getEnvInt("DB_MAX_IDLE", 10),
		},
		Keycloak: KeycloakConfig{
			URL:           getEnv("KEYCLOAK_URL", ""),
			Realm:         getEnv("KEYCLOAK_REALM", ""),
			ClientID:      getEnv("KEYCLOAK_CLIENT_ID", ""),
			ClientSecret:  getEnv("KEYCLOAK_CLIENT_SECRET", ""),
			AdminUser:     getEnv("KEYCLOAK_ADMIN_USER", ""),
			AdminPassword: getEnv("KEYCLOAK_ADMIN_PASSWORD", ""),
		},
		S3: map[string]S3Config{
			"public": {
				Endpoint:  getEnv("S3_PUBLIC_ENDPOINT", ""),
				Region:    getEnv("S3_PUBLIC_REGION", ""),
				AccessKey: getEnv("S3_PUBLIC_ACCESS_KEY", ""),
				SecretKey: getEnv("S3_PUBLIC_SECRET_KEY", ""),
				Bucket:    getEnv("S3_PUBLIC_BUCKET", ""),
				UseSSL:    getEnvBool("S3_PUBLIC_USE_SSL", false),
			},
			// "private": {
			// 	Endpoint:  getEnv("S3_PRIVATE_ENDPOINT", ""),
			// 	Region:    getEnv("S3_PRIVATE_REGION", ""),
			// 	AccessKey: getEnv("S3_PRIVATE_ACCESS_KEY", ""),
			// 	SecretKey: getEnv("S3_PRIVATE_SECRET_KEY", ""),
			// 	Bucket:    getEnv("S3_PRIVATE_BUCKET", ""),
			// 	UseSSL:    getEnvBool("S3_PRIVATE_USE_SSL", false),
			// },
		},
		ElasticApm: ElasticApmConfig{
			ServerURL:        getEnv("ELASTIC_APM_SERVER_URL", ""),
			ServiceName:      getEnv("ELASTIC_APM_SERVICE_NAME", "go-app-service"),
			ServiceVersion:   getEnv("ELASTIC_APM_SERVICE_VERSION", "1.0.0"),
			Environment:      getEnv("ELASTIC_APM_ENVIRONMENT", "development"),
			SecretToken:      getEnv("ELASTIC_APM_SECRET_TOKEN", ""),
			VerifyServerCert: getEnvBool("ELASTIC_APM_VERIFY_SERVER_CERT", false),
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

func getEnvBool(key string, defaultVal bool) bool {
	if val := viper.GetBool(key); val != false {
		return val
	}
	return defaultVal
}

func validate(cfg *Config) {
	if cfg.DB.Name == "" {
		log.Fatal("DB_NAME is required")
	}
}

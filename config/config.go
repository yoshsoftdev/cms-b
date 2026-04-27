package config

import (
	"os"
)

type Config struct {
	AppPort    string
	DBDriver   string
	DBDSN      string
	JWTSecret  string
	AdminUser  string
	AdminPass  string
	StorageURL string
}

func Load() Config {
	return Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBDriver:   getEnv("DB_DRIVER", "sqlite"),
		DBDSN:      getEnv("DB_DSN", "cms.db"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-change-me"),
		AdminUser:  getEnv("ADMIN_USER", "admin"),
		AdminPass:  getEnv("ADMIN_PASS", "admin123"),
		StorageURL: getEnv("STORAGE_BASE_URL", "https://cdn.example.com/uploads"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

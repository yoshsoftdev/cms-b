package config

import (
	"os"
)

type Config struct {
	AppPort          string
	DBDriver         string
	DBDSN            string
	JWTSecret        string
	AdminUser        string
	AdminPass        string
	StorageURL       string
	StorageDir       string
	DeployScriptPath string
}

func Load() Config {
	return Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		DBDriver:         getEnv("DB_DRIVER", "sqlite"),
		DBDSN:            getEnv("DB_DSN", "cms.db"),
		JWTSecret:        getEnv("JWT_SECRET", "super-secret-change-me"),
		AdminUser:        getEnv("ADMIN_USER", "admin"),
		AdminPass:        getEnv("ADMIN_PASS", "admin123"),
		StorageURL:       getEnv("STORAGE_BASE_URL", "/static/uploads"),
		StorageDir:       getEnv("STORAGE_DIR", "static/uploads"),
		DeployScriptPath: getEnv("DEPLOY_SCRIPT_PATH", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

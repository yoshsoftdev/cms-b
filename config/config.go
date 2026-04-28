package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
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
	loadEnvFiles(".env", ".env.local")
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

func loadEnvFiles(names ...string) {
	for _, name := range names {
		loadEnvFile(name)
	}
}

func loadEnvFile(name string) {
	abs, err := filepath.Abs(name)
	if err != nil {
		return
	}
	f, err := os.Open(abs)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

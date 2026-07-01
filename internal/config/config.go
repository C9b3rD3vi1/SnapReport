package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port            string
	DatabasePath    string
	UploadDir       string
	MaxUploadMB     int64
	MaxRequestBodyMB int64
	PDFDir          string
	Environment     string
	AllowedOrigins  []string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		DatabasePath:      getEnv("DATABASE_PATH", "./data/snapreport.db"),
		UploadDir:         getEnv("UPLOAD_DIR", "./uploads"),
		PDFDir:            getEnv("PDF_DIR", "./generated"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		MaxUploadMB:       getEnvInt("MAX_UPLOAD_MB", 10),
		MaxRequestBodyMB:  getEnvInt("MAX_REQUEST_BODY_MB", 50),
		AllowedOrigins:    getAllowedOrigins(getEnv("ALLOWED_ORIGINS", "*")),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return n
		}
	}
	return fallback
}

func getAllowedOrigins(val string) []string {
	if val == "*" {
		return []string{"*"}
	}
	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

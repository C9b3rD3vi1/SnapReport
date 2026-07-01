package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DatabasePath string
	UploadDir    string
	MaxUploadMB  int64
	PDFDir       string
	Environment  string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabasePath: getEnv("DATABASE_PATH", "./data/snapreport.db"),
		UploadDir:    getEnv("UPLOAD_DIR", "./uploads"),
		PDFDir:       getEnv("PDF_DIR", "./generated"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		MaxUploadMB:  getEnvInt("MAX_UPLOAD_MB", 10),
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

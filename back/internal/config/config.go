package config

import (
	"os"
	"time"
)

type Config struct {
	ServerPort    string
	DatabaseURL   string
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
	UploadDir     string
	MaxUploadSize int64
	Debug         bool
}

func Load() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://inventory_user:secret@localhost:5432/inventory?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "changeme"),
		JWTAccessTTL:  parseDuration(getEnv("JWT_ACCESS_TTL", "15m")),
		JWTRefreshTTL: parseDuration(getEnv("JWT_REFRESH_TTL", "168h")),
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		MaxUploadSize: int64(5 * 1024 * 1024), // 5 MB
		Debug:         getEnv("DEBUG", "false") == "true",
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}

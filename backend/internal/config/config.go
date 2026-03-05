package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	EncryptionKey string
	BLSAPIKey     string
}

func Load() (*Config, error) {
	godotenv.Load()

	encKey := getEnv("ENCRYPTION_KEY", "")
	if encKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY environment variable is required (32-byte hex string)")
	}

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/kroberts_personal?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
		EncryptionKey: encKey,
		BLSAPIKey:     getEnv("BLS_API_KEY", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

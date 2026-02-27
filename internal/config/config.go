package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	ServerPort  string
	Environment string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret string

	// QR Signing
	HMACSecret string

	// API
	APITimeout time.Duration
}

// Load reads configuration from environment variables
// It loads from .env.local (local dev), .env.dev (staging), or .env.prod (production)
func Load() (*Config, error) {
	// Determine which .env file to load
	// Priority: .env.local > .env.dev > .env.prod > .env (fallback)
	envFiles := []string{".env.local", ".env.dev", ".env.prod", ".env"}

	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			fmt.Printf("✓ Loaded configuration from %s\n", envFile)
			break
		}
	}
	// Continue even if no .env file found (env vars might be set)

	cfg := &Config{
		ServerPort:  getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		HMACSecret:  getEnv("HMAC_SECRET", "your-hmac-secret-change-in-production"),
		APITimeout:  30 * time.Second,
	}

	// Validate required fields
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	if cfg.Environment == "production" {
		if cfg.JWTSecret == "your-secret-key-change-in-production" {
			return nil, fmt.Errorf("JWT_SECRET must be set in production")
		}
		if cfg.HMACSecret == "your-hmac-secret-change-in-production" {
			return nil, fmt.Errorf("HMAC_SECRET must be set in production")
		}
	}

	return cfg, nil
}

// getEnv returns an environment variable or a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

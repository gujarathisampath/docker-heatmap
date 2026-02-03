package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	Port        string
	Environment string

	// Database
	DatabaseURL string

	// GitHub OAuth
	GitHubClientID     string
	GitHubClientSecret string
	GitHubCallbackURL  string

	// JWT
	JWTSecret string

	// Encryption
	EncryptionKey string

	// Frontend
	FrontendURL string

	// Docker Hub
	DockerHubAPIURL string
}

var AppConfig *Config

func Load() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		// Server
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/docker_heatmap?sslmode=disable"),

		// GitHub OAuth
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubCallbackURL:  getEnv("GITHUB_CALLBACK_URL", "http://localhost:8080/api/auth/github/callback"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),

		// Encryption (must be 32 bytes for AES-256)
		EncryptionKey: getEnv("ENCRYPTION_KEY", "a-32-byte-encryption-key-here!!"),

		// Frontend
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),

		// Docker Hub
		DockerHubAPIURL: getEnv("DOCKER_HUB_API_URL", "https://hub.docker.com/v2"),
	}

	// Validate required config
	if AppConfig.GitHubClientID == "" || AppConfig.GitHubClientSecret == "" {
		log.Println("Warning: GitHub OAuth credentials not configured")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

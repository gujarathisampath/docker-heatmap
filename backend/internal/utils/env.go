package utils

import (
	"os"
	"strconv"
)

// GetEnv reads an environment variable with a fallback default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt reads an environment variable as int with a fallback default value
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetEnvBool reads an environment variable as bool with a fallback default value
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

// IsProduction returns true if running in production environment
func IsProduction() bool {
	env := GetEnv("ENVIRONMENT", "development")
	return env == "production" || env == "prod"
}

// IsDevelopment returns true if running in development environment
func IsDevelopment() bool {
	return !IsProduction()
}

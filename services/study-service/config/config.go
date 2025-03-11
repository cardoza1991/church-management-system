package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the service
type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
	AuthService string // URL for the auth service for JWT verification
	ContactService string // URL for the contact service
}

// Load returns a new Config struct populated with values from environment variables
func Load() *Config {
	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "church_mgmt"),
		ServerPort:     getEnv("PORT", "8082"), // Different from other services
		AuthService:    getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
		ContactService: getEnv("CONTACT_SERVICE_URL", "http://localhost:8081"),
	}
}

// DSN returns a formatted database connection string
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
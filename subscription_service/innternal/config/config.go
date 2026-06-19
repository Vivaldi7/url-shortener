package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PGHost     string
	PGPort     string
	PGUser     string
	PGPassword string
	PGName     string
	SSLMode    string
	ServerPort string
	LogLevel   string
}

func LoadConfig() *Config {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Println("local.env file not found")
	}
	return &Config{
		PGHost:     getEnv("PG_Host", "localhost"),
		PGPort:     getEnv("PG_Port", "5432"),
		PGUser:     getEnv("PG_USER", "postgres"),
		PGPassword: getEnv("PG_PASSWORD", "postgres"),
		PGName:     getEnv("PG_NAME", "subscriptions"),
		SSLMode:    getEnv("SSL_MODE", "disable"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}
}

// Созданная фунция позволяет сократить код и не использовать os.Getenv так как он не проверяет на nil
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

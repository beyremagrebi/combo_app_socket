package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppHost      string
	AppPort      string
	ClientOrigin string
}

func Load() Config {
	loadEnvFile()

	return Config{
		AppHost:      getEnv("APP_HOST", "0.0.0.0"),
		AppPort:      getEnv("APP_PORT", "8800"),
		ClientOrigin: getEnv("CLIENT_ORIGIN", "http://localhost:5173"),
	}
}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
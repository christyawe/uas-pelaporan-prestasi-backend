// config/env.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}
}

func Env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
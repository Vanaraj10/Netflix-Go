package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DB_DSN    string
	JWTSecret string
	TMDBKey   string
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func mustGetEnv(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	dbDSN := mustGetEnv("DB_DSN")
	jwtSecret := mustGetEnv("JWT_SECRET")
	tmdbKey := mustGetEnv("TMDB_KEY")

	return &Config{
		Port:      port,
		DB_DSN:    dbDSN,
		JWTSecret: jwtSecret,
		TMDBKey:   tmdbKey,
	}
}

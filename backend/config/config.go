package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type config struct {
	Port       string
	dbDSN      string
	jwtSecret  string
	tmdbAPIKey string
	DB         *pgxpool.Pool
}

func LoadConfig() *config {
	_ = godotenv.Load()
	dbDSN := getEnv("DB_DSN", "")
	pool, err := pgxpool.New(context.Background(), dbDSN)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	return &config{
		Port:       getEnv("PORT", "8080"),
		dbDSN:      getEnv("DB_DSN", "user:password@tcp(localhost:3306)/dbname"),
		jwtSecret:  getEnv("JWT_SECRET", "defaultsecret"),
		tmdbAPIKey: getEnv("TMDB_API_KEY", "your_tmdb_api_key"),
		DB:         pool,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

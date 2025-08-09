package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Vanaraj10/Netflix/internal/config"
	"github.com/Vanaraj10/Netflix/internal/data"
	"github.com/Vanaraj10/Netflix/internal/handlers"
	"github.com/Vanaraj10/Netflix/internal/middleware"
	"github.com/go-chi/chi"
)

func main() {
	cfg := config.LoadConfig()

	db, err := data.NewDB(cfg.DB_DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Get("/v1/healthcheck", handlers.HealthcheckHandler)
	r.Post("/v1/users", handlers.RegisterHandler(db))
	r.Post("/v1/token", handlers.TokenHandler(db, cfg.JWTSecret))
	r.Get("/v1/discover/popular", handlers.PopularMoviesHandler(cfg.TMDBKey))
	r.Get("/v1/discover/search", handlers.SearchMoviesHandler(cfg.TMDBKey))

	// Protected watchlist routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(cfg.JWTSecret))
		r.Get("/v1/watchlist", handlers.WatchlistGetHandler(db, cfg.TMDBKey))
		r.Post("/v1/watchlist", handlers.WatchlistAddHandler(db))
		r.Patch("/v1/watchlist/{movie_id}", handlers.WatchlistUpdateHandler(db))
		r.Delete("/v1/watchlist/{movie_id}", handlers.WatchlistDeleteHandler(db))
	})

	fmt.Printf("Server will run on port: %s\n", cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

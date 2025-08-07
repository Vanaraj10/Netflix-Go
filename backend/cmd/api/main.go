package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Vanaraj10/Netflix/internal/config"
	"github.com/Vanaraj10/Netflix/internal/data"
	"github.com/Vanaraj10/Netflix/internal/handlers"
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

	fmt.Printf("Server will run on port: %s\n", cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

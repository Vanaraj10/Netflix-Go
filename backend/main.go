package main

import (
	"log"
	"net/http"

	"github.com/Vanaraj10/Netflix/config"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()
	defer cfg.DB.Close()

	r := chi.NewRouter()
	r.Get("/v1/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	log.Println("Starting server on :8080")
	http.ListenAndServe(":"+cfg.Port, r)
}

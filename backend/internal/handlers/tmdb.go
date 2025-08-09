package handlers

import (
	"io"
	"net/http"
	"net/url"
)

func PopularMoviesHandler(tmdbKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmdbUrl := "https://api.themoviedb.org/3/discover/movie?include_adult=false&include_video=false&language=en-US&page=1&sort_by=popularity.desc&api_key=d49063c816b83606d26ea2d89354a5d2"
		resp, err := http.Get(tmdbUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Failed to fetch popular movies", resp.StatusCode)
			return
		}
		// The response body is already written to w, so no need to do anything else
	}
}

func SearchMoviesHandler(tmdbKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		if query == "" {
			http.Error(w, "Query parameter is required", http.StatusBadRequest)
			return
		}
		tmdbURL := "https://api.themoviedb.org/3/search/movie?api_key=" + tmdbKey + "&query=" + url.QueryEscape(query)
		resp, err := http.Get(tmdbURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
	}
}

package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Vanaraj10/Netflix/internal/data"
	"github.com/Vanaraj10/Netflix/internal/middleware"
	"github.com/go-chi/chi"
)

type tmdbMovie struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	PosterPath  string `json:"poster_path"`
	Overview    string `json:"overview"`
	ReleaseDate string `json:"release_date"`
}

type watchlistItem struct {
	Movie      tmdbMovie `json:"movie"`
	Status     string    `json:"status"`
	UserRating int       `json:"user_rating"`
}

// --- GET /v1/watchlist ---

func WatchlistGetHandler(db *data.DB, tmdbKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		userMovies, err := db.GetUserMovies(ctx, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var wg sync.WaitGroup
		mu := sync.Mutex{}
		items := make([]watchlistItem, 0, len(userMovies))

		for _, um := range userMovies {
			wg.Add(1)
			go func(um data.UserMovie) {
				defer wg.Done()
				tmdbURL := "https://api.themoviedb.org/3/movie/" + strconv.FormatInt(um.MovieID, 10) + "?api_key=" + tmdbKey
				resp, err := http.Get(tmdbURL)
				if err != nil {
					return
				}
				defer resp.Body.Close()
				var movie tmdbMovie
				if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
					return
				}
				mu.Lock()
				items = append(items, watchlistItem{
					Movie:      movie,
					Status:     um.Status,
					UserRating: um.UserRating,
				})
				mu.Unlock()
			}(um)
		}
		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"watchlist": items,
		})
	}
}

// --- POST /v1/watchlist ---

type addWatchlistRequest struct {
	MovieID    int64  `json:"movie_id"`
	Status     string `json:"status"`
	UserRating int    `json:"user_rating"`
}

func WatchlistAddHandler(db *data.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		var req addWatchlistRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		um, err := db.AddUserMovie(ctx, userID, req.MovieID, req.Status, req.UserRating)
		if err != nil {
			http.Error(w, "Failed to add to watchlist", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"watchlist_item": um,
		})
	}
}

// --- PATCH /v1/watchlist/{movie_id} ---

type updateWatchlistRequest struct {
	Status     string `json:"status"`
	UserRating int    `json:"user_rating"`
}

func WatchlistUpdateHandler(db *data.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		movieIDStr := chi.URLParam(r, "movie_id")
		movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid movie_id", http.StatusBadRequest)
			return
		}
		var req updateWatchlistRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		if err := db.UpdateUserMovie(ctx, userID, movieID, req.Status, req.UserRating); err != nil {
			http.Error(w, "Failed to update watchlist", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// --- DELETE /v1/watchlist/{movie_id} ---

func WatchlistDeleteHandler(db *data.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		movieIDStr := chi.URLParam(r, "movie_id")
		movieID, err := strconv.ParseInt(movieIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid movie_id", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		if err := db.DeleteUserMovie(ctx, userID, movieID); err != nil {
			http.Error(w, "Failed to delete from watchlist", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

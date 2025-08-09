package data

import (
	"context"
	"time"
)

type UserMovie struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	MovieID    int64     `json:"movie_id"`
	Status     string    `json:"status"`
	UserRating int       `json:"user_rating"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AddUserMovie inserts a new movie into the user's watchlist.
func (db *DB) AddUserMovie(ctx context.Context, userID int64, movieID int64, status string, rating int) (*UserMovie, error) {
	var um UserMovie
	query := `
		INSERT INTO user_movies (user_id, movie_id, status, user_rating)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, movie_id, status, user_rating, created_at, updated_at
	`
	err := db.Pool.QueryRow(ctx, query, userID, movieID, status, rating).Scan(
		&um.ID, &um.UserID, &um.MovieID, &um.Status, &um.UserRating, &um.CreatedAt, &um.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &um, nil
}

// GetUserMovies retrieves all movies in the user's watchlist.
func (db *DB) GetUserMovies(ctx context.Context, userID int64) ([]UserMovie, error) {
	query := `
		SELECT id, user_id, movie_id, status, user_rating, created_at, updated_at
		FROM user_movies
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []UserMovie
	for rows.Next() {
		var um UserMovie
		if err := rows.Scan(&um.ID, &um.UserID, &um.MovieID, &um.Status, &um.UserRating, &um.CreatedAt, &um.UpdatedAt); err != nil {
			return nil, err
		}
		movies = append(movies, um)
	}
	return movies, nil
}

// UpdateUserMovie updates the status and rating for a movie in the user's watchlist.
func (db *DB) UpdateUserMovie(ctx context.Context, userID int64, movieID int64, status string, rating int) error {
	query := `
		UPDATE user_movies
		SET status = $1, user_rating = $2, updated_at = NOW()
		WHERE user_id = $3 AND movie_id = $4
	`
	_, err := db.Pool.Exec(ctx, query, status, rating, userID, movieID)
	return err
}

// DeleteUserMovie removes a movie from the user's watchlist.
func (db *DB) DeleteUserMovie(ctx context.Context, userID int64, movieID int64) error {
	query := `DELETE FROM user_movies WHERE user_id = $1 AND movie_id = $2`
	_, err := db.Pool.Exec(ctx, query, userID, movieID)
	return err
}

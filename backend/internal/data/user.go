package data

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

func (db *DB) CreateUser(ctx context.Context, name, email, password string) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	var user User
	query := `
		INSERT INTO users (name, email, hashed_password)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, hashed_password, created_at
	`
	err = db.Pool.QueryRow(ctx, query, name, email, string(hashed)).Scan(
		&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT id, name, email, hashed_password, created_at FROM users WHERE email = $1`
	err := db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

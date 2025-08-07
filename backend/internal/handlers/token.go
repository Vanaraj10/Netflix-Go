package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Vanaraj10/Netflix/internal/auth"
	"github.com/Vanaraj10/Netflix/internal/data"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TokenHandler(db *data.DB, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		req.Email = strings.TrimSpace(req.Email)
		req.Password = strings.TrimSpace(req.Password)
		if req.Email == "" || req.Password == "" {
			http.Error(w, "Email and Password are required", http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user, err := db.GetUserByEmail(ctx, req.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(req.Password)); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := auth.GenerateJWT(user.ID, secret)
		if err != nil {
			http.Error(w, "Error generating Auth Token", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}

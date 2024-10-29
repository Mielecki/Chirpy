package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mielecki/Chirpy/internal/auth"
	"github.com/Mielecki/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string `json:"email"`
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)

	data := parameters{}
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding error", err)
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Hashing password error", err)
	} 

	userData, err := cfg.database.CreateUser(req.Context(), database.CreateUserParams{
		Email: data.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Creating user error", err)
		return
	}
	
	respondWithJSON(w, 201, User{
		ID: userData.ID,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		Email: userData.Email,
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}

	type returnVals struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(req.Body)

	data := parameters{}
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding error", err)
		return
	}

	user, err := cfg.database.GetUserByEmail(req.Context(), data.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Getting user error", err)
		return
	}

	if err := auth.CheckPasswordHash(data.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return 
	}

	expiresInSeconds := data.ExpiresInSeconds
	if expiresInSeconds == 0 || expiresInSeconds > 3600{
		expiresInSeconds = 3600
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expiresInSeconds) * time.Second)
	if err != nil {
		respondWithError(w, 500, "Making token error", err)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
		Token: token,
	})
}
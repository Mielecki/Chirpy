package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Mielecki/Chirpy/internal/auth"
	"github.com/Mielecki/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Getting token error", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	data := parameters{}
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding error", err)
		return
	}

	if ok := validateChirp(&data.Body); !ok {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.database.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: data.Body,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Creating chrip error", err)
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func validateChirp(body *string) bool {
	if len(*body) > 140 {
		return false
	}

	replaceProfane(body)
	return true
}

func replaceProfane(s *string) {
	splitted := strings.Split(*s, " ")

	for i, item := range splitted {
		item = strings.ToLower(item)
		if item == "kerfuffle" || item == "sharbert" || item == "fornax"{
			splitted[i] = "****"
		}
	}

	*s = strings.Join(splitted, " ")
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.database.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Getting chrips error", err)
		return
	}

	chirpsJSON := []Chirp{}
	
	for _, chirp := range chirps {
		chirpsJSON = append(chirpsJSON, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, 200, chirpsJSON)
}


func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Parsing chirpID error", err)
		return
	}

	chirp, err := cfg.database.GetChirp(req.Context(), chirpID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, 404, "No chirp error", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Getting chirps error", err)
		return
	}

	respondWithJSON(w, 200, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}
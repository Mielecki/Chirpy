package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Mielecki/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolka(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	if key, err := auth.GetAPIKey(req.Header); err != nil || key != cfg.polkaKey{
		respondWithError(w, http.StatusUnauthorized, "Invalid Polka key", err)
		return
	}


	data := parameters{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding error", err)
		return
	}

	if data.Event != "user.upgraded" {
		respondWithError(w, 204, "Unknown event", errors.New("unknown event type"))
		return
	}

	if _, err := cfg.database.UpgradeToChripyRed(req.Context(), data.Data.UserID); err != nil {
		respondWithError(w, 404, "upgrading error", err)
		return
	}

	respondWithJSON(w, 204, struct{}{}) 
}
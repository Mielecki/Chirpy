package main

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/Mielecki/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database *database.Queries
	platform string
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(403)
		w.Write([]byte("Forbidden"))
	}
	if err := cfg.database.Reset(context.Background()); err != nil {
		respondWithError(w, 500, "Resetting error", err)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
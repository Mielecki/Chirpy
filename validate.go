package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	data := parameters{}
	if err := decoder.Decode(&data); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Decoding error", err)
		return
	}

	if len(data.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	replaceProfane(&data.Body)
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: data.Body,
	})
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
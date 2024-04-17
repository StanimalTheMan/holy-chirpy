package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/StanimalTheMan/holy-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersUpgrade(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		UserID int `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	err := auth.ValidateAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid api key")
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	// upgrade user status to red status
	err = cfg.DB.UpgradeUser(params.Data.UserID)

	if errors.Is(err, os.ErrNotExist) {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, "User upgraded to chirpy red status")
}

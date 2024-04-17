package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/StanimalTheMan/holy-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// authenticated endpoint
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	// only allow deletion if user is author of chirp
	subject, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID, userIDInt)
	if errors.Is(err, os.ErrNotExist) {
		respondWithError(w, http.StatusNotFound, "Chirp does not exist")
		return
	}
	if err != nil {
		// user is not author of chirp
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	respondWithJSON(w, 200, "Chirp successfully deleted")
}

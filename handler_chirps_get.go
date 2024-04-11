package main

import (
	"errors"
	"net/http"
	"sort"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirp, err := cfg.DB.GetChirp(r.PathValue("chirpID"))
	if err != nil {
		if err.Error() == errors.New("chirp does not exist").Error() {
			respondWithError(w, http.StatusNotFound, "Chirp does not exist")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirp")
		return
	}
	respondWithJSON(w, http.StatusOK, dbChirp)
}

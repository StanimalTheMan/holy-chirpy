package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/StanimalTheMan/holy-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		// ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	defaultAccessTokenExpiration := 60 * 60
	defaultRefreshTokenExpiration := 60 * 60 * 24 * 60
	// if params.ExpiresInSeconds == 0 {
	// 	params.ExpiresInSeconds = defaultExpiration
	// } else if params.ExpiresInSeconds > defaultExpiration {
	// 	params.ExpiresInSeconds = defaultExpiration
	// }

	access_token, err := auth.MakeJWT(user.ID, "chirpy-access", cfg.jwtSecret, time.Duration(defaultAccessTokenExpiration)*time.Second)
	refresh_token, err := auth.MakeJWT(user.ID, "chirpy-refresh", cfg.jwtSecret, time.Duration(defaultRefreshTokenExpiration)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token:        access_token,
		RefreshToken: refresh_token,
	})
}

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/StanimalTheMan/holy-chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		AccessToken string `json:"token"`
	}
	tokenString, err := auth.GetBearerToken(r.Header)
	fmt.Printf("Trying to refresh with %v", tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(cfg.jwtSecret), nil },
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Reject requests if access token in header is not a refresh token (check issuer).
	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	// check if refresh token is revoked
	isRefreshTokenRevoked := cfg.DB.CheckIfRevoked(tokenString)
	if isRefreshTokenRevoked {
		fmt.Println("SHOULD BE HERE")
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	userID, _ := strconv.Atoi(userIDString)
	fmt.Println("put attempt", userID, issuer, cfg.jwtSecret, 60*60*time.Second)
	newAccessToken, _ := auth.MakeJWT(userID, "chirpy-access", cfg.jwtSecret, 60*60*time.Second)
	fmt.Println("new ACCESS", newAccessToken)
	respondWithJSON(w, http.StatusOK, response{
		AccessToken: newAccessToken,
	})
}

func (cfg *apiConfig) HandleRevokeToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Println("error0")
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(cfg.jwtSecret), nil },
	)
	if err != nil {
		fmt.Println("error1")
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		fmt.Println("error2")
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Reject requests if access token in header is a refresh token (check issuer).
	if issuer == "chirpy-access" {
		fmt.Println("error3")
		respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	cfg.DB.CreateRevokedToken(tokenString)
	respondWithJSON(w, http.StatusOK, "successfully revoked")
}

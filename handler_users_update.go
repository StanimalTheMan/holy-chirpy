package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/StanimalTheMan/holy-chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type CustomClaims struct {
	jwt.RegisteredClaims
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	// Decode JSON Request Body
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		Email string `json:"email"`
		ID    int    `json:"id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// TODO: probably should validate email later

	// Validate JWT
	// Extract token from request headers
	fmt.Println(r.Header.Get("Authorization"))
	tokenString := strings.Split(r.Header.Get("Authorization"), " ")[1]
	fmt.Println("Extracted token:", tokenString)

	// Load environment variables from the .env file
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Get value of JWT Secret from .env file
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorized")
		return
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
		return
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error converting user ID")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, _ := cfg.DB.UpdateUser(userId, params.Email, hashedPassword)

	respondWithJSON(w, http.StatusOK, response{
		ID:    user.ID,
		Email: user.Email,
	})
}

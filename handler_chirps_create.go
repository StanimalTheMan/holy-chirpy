package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/StanimalTheMan/holy-chirpy/internal/auth"
	"github.com/joho/godotenv"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	// Decode JSON Request Body
	type parameters struct {
		Body string `json:"body"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
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

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	url := "https://api.openai.com/v1/moderations"

	type InputData struct {
		Input string `json:"input"`
	}
	params = struct {
		Body string `json:"body"`
	}{
		Body: params.Body,
	}

	inputData := InputData{
		Input: params.Body,
	}

	jsonStr, err := json.Marshal(inputData)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating request", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	// Load environment variables from the .env file
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Get value of "OPENAI_API_KEY" environment variable from .env file
	apiKey := os.Getenv("OPENAI_API_KEY")

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp, err)

	// body := make([]byte, 1024)
	// _, err = resp.Body.Read(body)
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return
	// }
	// fmt.Println("Response Body:", string(body))

	var response struct {
		ID      string `json:"id"`
		Model   string `json:"model"`
		Results []struct {
			Flagged bool `json:"flagged"`
		} `json:"results"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error reading API response")
		return
	}

	// Check if the content is flagged by OpenAI API
	if len(response.Results) > 0 && response.Results[0].Flagged {
		respondWithError(w, http.StatusBadRequest, "Content flagged by OpenAI")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned, userIDInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:       chirp.ID,
		Body:     chirp.Body,
		AuthorID: userIDInt,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	cleaned := getCleanedBody(body)
	return cleaned, nil
}

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// OpenAI moderation endpoint
	url := "https://api.openai.com/v1/moderations"

	jsonStr := []byte(`{
		"input": "You have a big face.  Get plastic surgery"
	}`)

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

	fmt.Println("Response Status:", resp.Status)

	body := make([]byte, 1024)
	_, err = resp.Body.Read(body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(body))
}

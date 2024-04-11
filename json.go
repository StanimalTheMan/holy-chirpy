package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func postTweetHandler(w http.ResponseWriter, r *http.Request) {
	// Decode JSON Request Body
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	// Encode JSON Response Body
	type errorReturnVals struct {
		Error string `json:"string"`
	}
	if len(params.Body) > 140 {
		respBody := errorReturnVals{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			respBody.Error = fmt.Sprintf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(data)
	} else {
		type successReturnVals struct {
			Valid bool `json:"valid"`
		}
		respBody := successReturnVals{
			Valid: true,
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			respBody := errorReturnVals{
				Error: fmt.Sprintf("Error marshalling JSON: %s", err),
			}
			w.WriteHeader(500)
			w.Write([]byte(respBody.Error))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(data)
	}
}

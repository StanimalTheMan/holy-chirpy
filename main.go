package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/StanimalTheMan/holy-chirpy/database"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(fileServerHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", handlerChirpsValidate)
	mux.HandleFunc("GET /api/chirps", handleFetchChirps)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	// Decode JSON Request Body
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := getCleanedBody(params.Body)

	// respondWithJSON(w, http.StatusOK, returnVals{
	// 	CleanedBody: cleaned,
	// })
	db, err := database.NewDB("database.json")

	if err != nil {
		fmt.Println("Err1")
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	chirp, err := db.CreateChirp(cleaned)

	if err != nil {
		fmt.Println("Err2")
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}
	fmt.Println("chirp", chirp)
	respondWithJSON(w, http.StatusCreated, chirp)
}

func handleFetchChirps(w http.ResponseWriter, r *http.Request) {
	db, err := database.NewDB("database.json")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Server")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Encode JSON response body
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

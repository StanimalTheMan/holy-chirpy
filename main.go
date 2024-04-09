package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Create a new http.ServeMux
	mux := http.NewServeMux()

	// Wrap mux in custom middleware that adds CORS headers to response
	corsMux := middlewareCors(mux)

	// Create a new http.Server and use corsMux on handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: corsMux,
	}

	// Start HTTP server
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

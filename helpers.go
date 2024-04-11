package main

import (
	"strings"
)

func getCleanedBody(msg string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	splitMsg := strings.Split(msg, " ")
	for i, word := range splitMsg {
		normalizedWord := strings.ToLower(word)
		if _, ok := badWords[normalizedWord]; ok {
			splitMsg[i] = "****"
		}
	}
	return strings.Join(splitMsg, " ")
}

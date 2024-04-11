package main

import (
	"fmt"
	"strings"
)

func replaceBadWords(msg string) string {
	cleanedMsg := []string{}
	splitMsg := strings.Split(msg, " ")
	for _, word := range splitMsg {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			cleanedMsg = append(cleanedMsg, "****")
		} else {
			cleanedMsg = append(cleanedMsg, word)
		}
	}
	fmt.Println(cleanedMsg)
	return strings.Join(cleanedMsg, " ")
}

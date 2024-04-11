package main

import (
	"testing"
)

func TestReplaceBadWords(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"You Sharbert you mother kerfuffle", "You **** you mother ****"},
		{"Kerfuffle is a bad word and so is kerfuffle", "**** is a bad word and so is ****"},
		{"Your mouth is holy.  No bad words here.", "Your mouth is holy.  No bad words here."},
		{"Haha you kerfuffle.  I can curse with . or ! ya sharbert.", "Haha you kerfuffle.  I can curse with . or ! ya sharbert."},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := replaceBadWords(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

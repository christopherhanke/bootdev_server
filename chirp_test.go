package main

import (
	"testing"
)

func TestReplaceBadWords(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"This contains kerfuffle", "This contains ****"},
		{"No bad words here", "No bad words here"},
		{"sharbert is here", "**** is here"},
		{"Fornax", "****"},
		{"sharbert!", "sharbert!"},
		{"KERFUFFLE", "****"},
	}

	for _, test := range tests {
		result := replaceBadWords(test.input)
		if result != test.expected {
			t.Errorf("Excpected '%s', but got '%s'", test.expected, result)
		}
	}
}

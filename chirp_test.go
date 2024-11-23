package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestHandlerValidate(t *testing.T) {
	reqBody := `{"body": "This is a kerfuffle and a Sharbert"}`
	req := httptest.NewRequest("POST", "/api/validate_chirp", strings.NewReader(reqBody))
	rec := httptest.NewRecorder()
	handlerValidate(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var response map[string]string
	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	_, exists := response["cleaned_body"]
	if !exists {
		t.Fatalf("Expected cleaned_body key to be present in response.")
	}

	expected := "This is a **** and a ****"
	if response["cleaned_body"] != expected {
		t.Errorf("Expected cleaned_body to be '%s', but got '%s'", expected, response["cleaned_body"])
	}
}

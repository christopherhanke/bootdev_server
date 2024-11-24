package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

// handler which validates incoming chirps
func handlerValidate(respw http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(respw, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp ist too long: %v", len(params.Body))
		respondWithError(respw, http.StatusBadRequest, "Chirp is too long")
		return
	}

	log.Printf("valid chirp")
	response := map[string]string{"cleaned_body": replaceBadWords(params.Body)}
	respondWithJSON(respw, http.StatusOK, response)
}

// replacing all bad words in a string with ****
func replaceBadWords(msg string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(msg, " ")
	for key, word := range words {
		for _, badWord := range badWords {
			if badWord == strings.ToLower(word) {
				words[key] = "****"
				continue
			}
		}
	}
	return strings.Join(words, " ")
}

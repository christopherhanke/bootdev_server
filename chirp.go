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

type returnMsg struct {
	Error       string `json:"error"`
	CleanedBody string `json:"cleaned_body"`
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

// response with error, giving http-StatusCode and error message
func respondWithError(respw http.ResponseWriter, code int, msg string) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(code)
	respmsg, err := json.Marshal(returnMsg{Error: msg})
	if err != nil {
		http.Error(respw, err.Error(), http.StatusInternalServerError)
		return
	}
	respw.Write(respmsg)
}

// response with JSON data, giving http-StatusCode and payload for response
func respondWithJSON(respw http.ResponseWriter, code int, payload interface{}) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(code)
	respmsg, err := json.Marshal(payload)
	if err != nil {
		http.Error(respw, err.Error(), http.StatusInternalServerError)
		return
	}
	respw.Write(respmsg)
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

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/christopherhanke/bootdev_server/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// handler which validates incoming chirps
func (cfg *apiConfig) handlerChirp(respw http.ResponseWriter, req *http.Request) {
	// format of the incoming json data
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// decode incoming request to handle
	decoder := json.NewDecoder(req.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(respw, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// message body has to be max 140 characters
	if len(params.Body) > 140 {
		log.Printf("Chirp ist too long: %v", len(params.Body))
		respondWithError(respw, http.StatusBadRequest, "Chirp is too long")
		return
	}
	log.Printf("valid chirp")

	// create chirp in database and handle error
	chirp, err := cfg.databaseQueries.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   replaceBadWords(params.Body),
		UserID: params.UserID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(respw, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// handle chirp to marshal for json
	response := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(respw, http.StatusCreated, response)
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

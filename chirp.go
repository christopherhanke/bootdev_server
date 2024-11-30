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
	log.Printf("valid chirp: %s", params.Body)

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

// handler gets all chirps in ascending order from database and parses json
func (cfg *apiConfig) handlerGetChirps(respw http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.databaseQueries.GetChirps(req.Context())
	if err != nil {
		log.Printf("Error loading chirps from database: %s", err)
		respondWithError(respw, http.StatusInternalServerError, "error loading chirps")
		return
	}

	var jsonChirps []Chirp
	for _, chirp := range chirps {
		val := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		jsonChirps = append(jsonChirps, val)
	}

	log.Printf("request all chirps: %v", len(jsonChirps))
	respondWithJSON(respw, http.StatusOK, jsonChirps)
}

func (cfg *apiConfig) handlerGetChip(respw http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")
	if id == "" {
		log.Printf("no chirp found: %s", req.URL.Path)
		respondWithError(respw, http.StatusNotFound, "no chirp found")
		return
	}

	val, err := uuid.Parse(id)
	if err != nil {
		log.Printf("error reading uuid: %s", id)
		respondWithError(respw, http.StatusNotFound, "could not read uuid")
		return
	}

	chirp, err := cfg.databaseQueries.GetChirp(req.Context(), val)
	if err != nil {
		log.Printf("no chirp found: %v", err)
		respondWithError(respw, http.StatusNotFound, "error searching chirp")
	}
	jsonChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	log.Printf("chirp ID: %s", chirp.ID)
	respondWithJSON(respw, http.StatusOK, jsonChirp)
}

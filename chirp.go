package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/christopherhanke/bootdev_server/internal/auth"
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
		Body string `json:"body"`
	}

	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("get BearerToken failed: %s", err)
		respondWithError(respw, http.StatusUnauthorized, "access denied")
		return
	}
	userID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		log.Printf("failed to validate token: %s", err)
		respondWithError(respw, http.StatusUnauthorized, "access denied")
		return
	}

	// decode incoming request to handle
	decoder := json.NewDecoder(req.Body)
	var params parameters
	err = decoder.Decode(&params)
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
		UserID: userID,
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

// handler gets all chirps in ascending order from database and parses json. it checks for the author ID in query.
func (cfg *apiConfig) handlerGetChirps(respw http.ResponseWriter, req *http.Request) {
	var chirps []database.Chirp
	var err error

	// check if URL query has authorID, when authorID is given, only return chirps from that author
	author := req.URL.Query().Get("author_id")
	if author != "" {
		userID, err := uuid.Parse(author)
		if err != nil {
			log.Printf("error parsing userID: %s", err)
			respondWithError(respw, http.StatusInternalServerError, "could not read author id")
			return
		}
		chirps, err = cfg.databaseQueries.GetChirpsAuthor(req.Context(), userID)
		if err != nil {
			log.Printf("error getting chirps for user: %s", err)
			respondWithError(respw, http.StatusInternalServerError, "error loading chirps")
			return
		}
	} else {
		chirps, err = cfg.databaseQueries.GetChirps(req.Context())
		if err != nil {
			log.Printf("Error loading chirps from database: %s", err)
			respondWithError(respw, http.StatusInternalServerError, "error loading chirps")
			return
		}
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

	if author != "" {
		log.Printf("request all chirps from %s: %v", author, len(jsonChirps))
	} else {
		log.Printf("request all chirps: %v", len(jsonChirps))
	}
	respondWithJSON(respw, http.StatusOK, jsonChirps)
}

func (cfg *apiConfig) handlerGetChip(respw http.ResponseWriter, req *http.Request) {
	// get ID from PathValue, parse it to UUID and check for chirp in database.
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
		return
	}

	// parse Chirp data to JSON aligned struct Chirp and make HTTP response
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

// delete chirp by given ID in PathValue
func (cfg *apiConfig) handlerDeleteChirp(respw http.ResponseWriter, req *http.Request) {
	// get ID from PathValue, parse it to UUID and check for chirp in database.
	chirpID := req.PathValue("chirpID")
	if chirpID == "" {
		log.Printf("chirp ID is missing: %s", req.URL.Path)
		respondWithError(respw, http.StatusNotFound, "chirp ID is missing")
		return
	}
	val, err := uuid.Parse(chirpID)
	if err != nil {
		log.Printf("error reading uuid: %s", chirpID)
		respondWithError(respw, http.StatusNotFound, "could not read uuid")
		return
	}
	chirp, err := cfg.databaseQueries.GetChirp(req.Context(), val)
	if err != nil {
		log.Printf("no chirp found: %s", err)
		respondWithError(respw, http.StatusNotFound, "error searching chirp")
		return
	}

	// check for authentication Token and get UserID from it. compare author of chirp with user authentication.
	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("get BearerToken failed: %s", err)
		respondWithError(respw, http.StatusUnauthorized, "access denied")
		return
	}
	userID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		log.Printf("failed to validate token: %s", err)
		respondWithError(respw, http.StatusUnauthorized, "access denied")
		return
	}
	if chirp.UserID != userID {
		log.Printf("user is not author of chirp: %s", chirp.UserID)
		respondWithError(respw, http.StatusForbidden, "deletion is not authorized")
		return
	}

	// delete chirp in database
	err = cfg.databaseQueries.DeleteChirp(req.Context(), val)
	if err != nil {
		log.Printf("failed to delete chirp: %s", err)
		respondWithError(respw, http.StatusNotFound, "error on deletion")
		return
	}

	respondWithJSON(respw, http.StatusNoContent, Chirp{})
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// handler to create new user
func (cfg *apiConfig) handlerCreateUser(respw http.ResponseWriter, req *http.Request) {
	// format of the incoming json data
	type parameters struct {
		Email string `json:"email"`
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

	user, err := cfg.databaseQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(respw, http.StatusInternalServerError, "couldn't create user")
		return
	}

	respuser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	log.Printf("user created: %s", respuser.Email)
	respondWithJSON(respw, http.StatusCreated, respuser)
}

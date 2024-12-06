package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/christopherhanke/bootdev_server/internal/auth"
	"github.com/christopherhanke/bootdev_server/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

// format of the incoming json data for user
type parameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// handler to create new user
func (cfg *apiConfig) handlerCreateUser(respw http.ResponseWriter, req *http.Request) {
	// decode incoming request to handle
	params, err := decodeIncomingUser(req)
	if err != nil {
		respondWithError(respw, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// error if password is invalid, currently empty #TODO
	if params.Password == "" {
		log.Printf("Passoword not valid: too short")
		respondWithError(respw, http.StatusForbidden, "Password invalid")
		return
	}

	// hash password for storing in database
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(respw, http.StatusInternalServerError, "Error hashing password")
		return
	}

	// create user for database
	user, err := cfg.databaseQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(respw, http.StatusInternalServerError, "couldn't create user")
		return
	}

	// create user for response without password or hash
	respuser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	log.Printf("user created: %s", respuser.Email)
	respondWithJSON(respw, http.StatusCreated, respuser)
}

// handler to login user
func (cfg *apiConfig) handlerLoginUser(respw http.ResponseWriter, req *http.Request) {
	// decode incoming request to handle
	params, err := decodeIncomingUser(req)
	if err != nil {
		respondWithError(respw, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// get user by email
	user, err := cfg.databaseQueries.GetUser(req.Context(), params.Email)
	if err != nil {
		log.Printf("Error searching for user: %s", params.Email)
		respondWithError(respw, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Invalid password: %s", params.Password)
		respondWithError(respw, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	// create user for response without password or hash
	respuser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	log.Printf("user logged in: %s", respuser.Email)
	respondWithJSON(respw, http.StatusOK, respuser)
}

// decode incoming user data to given struct
func decodeIncomingUser(req *http.Request) (parameters, error) {
	// decode incoming request to handle
	decoder := json.NewDecoder(req.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		return parameters{}, err
	}
	return params, nil
}

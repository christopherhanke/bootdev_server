package main

import (
	"log"
	"net/http"
	"time"

	"github.com/christopherhanke/bootdev_server/internal/auth"
)

type refreshToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(respw http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error reading bearer token: %s", err)
		respondWithError(respw, http.StatusUnauthorized, "invalid token")
		return
	}

	dbEntry, err := cfg.databaseQueries.GetUserFromRefreshToken(req.Context(), token)
	if err != nil || dbEntry.ExpiresAt.Before(time.Now()) {
		respondWithError(respw, http.StatusUnauthorized, "something went wrong")
		return
	}

	if dbEntry.RevokedAt.Valid {
		log.Printf("revoked: %s", dbEntry.RevokedAt.Time)
		// dbEntry.RevokedAt.Time.Before(time.Now())
		respondWithError(respw, http.StatusUnauthorized, "something went wrong")
		return

	}

	newToken, err := auth.MakeJWT(dbEntry.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(respw, http.StatusInternalServerError, "something went wrong")
		return
	}

	response := refreshToken{
		Token: newToken,
	}
	log.Print("refresh succesful")
	respondWithJSON(respw, http.StatusOK, response)
}

func (cfg *apiConfig) handlerRevoke(respw http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Printf("error reading bearer token: %s", err)
		respondWithError(respw, http.StatusUnauthorized, "invalid token")
		return
	}
	log.Printf("Token: %s", token)
	err = cfg.databaseQueries.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		log.Printf("error revoke token: %s", err)
		respondWithError(respw, http.StatusInternalServerError, "something went wrong")
		return
	}

	log.Printf("revoked token: %s", token)
	respw.WriteHeader(http.StatusNoContent)
}

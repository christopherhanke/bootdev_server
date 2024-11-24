package main

import (
	"log"
	"net/http"
)

// reset hits on HTTP server
func (cfg *apiConfig) handlerReset(respw http.ResponseWriter, req *http.Request) {
	// Reset is restricted to dev
	if cfg.enviroment != "dev" {
		respondWithError(respw, http.StatusForbidden, "Access Denied")
		return
	}

	// reset Hits to zero
	cfg.fileserverHits.Store(0)

	// reset user database
	err := cfg.databaseQueries.DeleteUsers(req.Context())
	if err != nil {
		log.Printf("Deleting users failed: %s", err)
		respondWithError(respw, http.StatusInternalServerError, "Deleting users failed")
	}

	respw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	respw.WriteHeader(http.StatusOK)
	respw.Write([]byte("Server reset"))
	log.Printf("Server reset")
}

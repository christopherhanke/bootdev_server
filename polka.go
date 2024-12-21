package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/christopherhanke/bootdev_server/internal/auth"
	"github.com/google/uuid"
)

// data format of incoming data
type webhooksParameters struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(respw http.ResponseWriter, req *http.Request) {
	// check for authorization
	key, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Printf("API key failed: %s", err)
		respw.WriteHeader(http.StatusUnauthorized)
		return
	}
	if key != cfg.polkaKey {
		log.Printf("unauthorized request on Polka Webhook: %s", req.Host)
		respw.WriteHeader(http.StatusUnauthorized)
		return
	}

	// decode incoming json on http request
	decoder := json.NewDecoder(req.Body)
	var params webhooksParameters
	err = decoder.Decode(&params)
	if err != nil {
		respw.WriteHeader(http.StatusBadRequest)
		return
	}

	// check for event data
	if params.Event != "user.upgraded" {
		log.Printf("no user upgrade: %s", params.Event)
		respw.WriteHeader(http.StatusNoContent)
		return
	}

	if params.Event == "user.upgraded" {
		userID, err := uuid.Parse(params.Data.UserID)
		if err != nil {
			log.Printf("parsing user id from event failed: %s", err)
			respw.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = cfg.databaseQueries.UpgradeUser(req.Context(), userID)
		if err != nil {
			log.Printf("could not upgrade user: %s", err)
			respw.WriteHeader(http.StatusNotFound)
			return
		}
		log.Printf("user upgraded: %s", userID)
		respw.WriteHeader(http.StatusNoContent)
	}

}

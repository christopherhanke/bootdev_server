package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// handler which validates incoming chirps
func handlerValidate(respw http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnMsg struct {
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respw.WriteHeader(http.StatusInternalServerError)
		msg, err := json.Marshal(returnMsg{Error: "Something went wrong"})
		if err == nil {
			respw.Header().Set("Content-Type", "application/json")
			respw.Write(msg)
		}
		return
	}

	if len(params.Body) > 140 {
		log.Printf("Chirp ist too long: %v", len(params.Body))
		respw.WriteHeader(http.StatusBadRequest)
		msg, err := json.Marshal(returnMsg{Error: "Chirp is too long"})
		if err == nil {
			respw.Header().Set("Content-Type", "application/json")
			respw.Write(msg)
		}
		return
	}

	log.Printf("valid chirp")
	respw.WriteHeader(http.StatusOK)
	msg, err := json.Marshal(returnMsg{Valid: true})
	if err == nil {
		respw.Header().Set("Content-Type", "application/json")
		respw.Write(msg)
	}
}

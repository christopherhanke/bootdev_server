package main

import (
	"encoding/json"
	"net/http"
)

type returnMsg struct {
	Error       string `json:"error"`
	CleanedBody string `json:"cleaned_body"`
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

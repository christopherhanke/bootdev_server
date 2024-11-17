package main

import (
	"net/http"
	"strconv"
)

// reset hits on HTTP server
func (cfg *apiConfig) handlerReset(respw http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	hits := "Hits: " + strconv.FormatInt(int64(cfg.fileserverHits.Load()), 10)
	respw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	respw.WriteHeader(200)
	respw.Write([]byte(hits))
}

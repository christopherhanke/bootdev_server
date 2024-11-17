package main

import (
	"net/http"
	"strconv"
)

// serve hits on HTTP Fileserver in plain text
func (cfg *apiConfig) handlerHits(respw http.ResponseWriter, req *http.Request) {
	hits := "Hits: " + strconv.FormatInt(int64(cfg.fileserverHits.Load()), 10)

	respw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	respw.WriteHeader(200)
	respw.Write([]byte(hits))
}

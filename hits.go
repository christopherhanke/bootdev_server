package main

import (
	"net/http"
	"strconv"
)

// serve hits on HTTP Fileserver in plain text
func (cfg *apiConfig) handlerHits(respw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		respw.WriteHeader(405)
		respw.Write([]byte("not get\n"))
		return
	}

	hits := "Hits: " + strconv.FormatInt(int64(cfg.fileserverHits.Load()), 10) + "\n"

	respw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	respw.WriteHeader(200)
	respw.Write([]byte(hits))
}

package main

import "net/http"

func hanlderHealthz(respw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		respw.WriteHeader(http.StatusMethodNotAllowed)
		respw.Write([]byte("not get\n"))
		return
	}
	respw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	respw.WriteHeader(http.StatusOK)
	respw.Write([]byte("OK\n"))
}

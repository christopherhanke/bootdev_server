package main

import "net/http"

func hanlderHealthz(respw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		respw.WriteHeader(405)
		respw.Write([]byte("not get"))
		return
	}
	respw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	respw.WriteHeader(200)
	respw.Write([]byte("OK"))
}

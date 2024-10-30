package main

import (
	"log"
	"net/http"
)

func main() {
	new := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: new,
	}
	new.Handle("/", http.FileServer(http.Dir(".")))
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed: %v\n", err)
	}

}

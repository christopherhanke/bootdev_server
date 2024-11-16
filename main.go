package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	new := http.NewServeMux()
	new.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	server := http.Server{
		Addr:    ":" + port,
		Handler: new,
	}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed: %v\n", err)
	}

}

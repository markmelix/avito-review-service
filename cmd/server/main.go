package main

import (
	"log"
	"net/http"
)

const httpServerAddress string = "0.0.0.0:8080"

func main() {
	log.SetFlags(log.LstdFlags)

	mux := http.NewServeMux()

	log.Printf("Starting http-server on %s\n", httpServerAddress)
	if err := http.ListenAndServe(httpServerAddress, mux); err != nil {
		log.Fatal("Error while serving http-server: %w", err)
	}
}

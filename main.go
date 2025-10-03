package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	fmt.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on :8080: %v\n", err)
	}
}
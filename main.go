package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	const filepathRoot = "."
	const port = "8080"

	fs := http.FileServer(http.Dir(filepathRoot))
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	
	fmt.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on :%s: %v\n", port, err)
	}
}
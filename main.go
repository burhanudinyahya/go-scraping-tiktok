package main

import (
	"go-tiktok-scraping/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/search", handlers.VideoSearch).Methods("GET")
	r.HandleFunc("/video", handlers.VideoDetail).Methods("GET")

	http.Handle("/", r)
	port := os.Getenv("PORT") // Render sets this automatically
	if port == "" {
		port = "10000"
	}

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

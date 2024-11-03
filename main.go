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
	r.HandleFunc("/search", handlers.Search).Methods("GET")
	r.HandleFunc("/video", handlers.VideoDetail).Methods("GET")
	r.HandleFunc("/stream", handlers.Stream).Methods("GET")

	http.Handle("/", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}

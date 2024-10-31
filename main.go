package main

import (
	"go-tiktok-scraping/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/search", handlers.VideoSearch).Methods("GET")
	r.HandleFunc("/video", handlers.VideoDetail).Methods("GET")

	http.Handle("/", r)
	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

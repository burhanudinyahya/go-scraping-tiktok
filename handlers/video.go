package handlers

import (
	"encoding/json"
	"fmt"
	"go-tiktok-scraping/services"
	"io"
	"net/http"
)

func VideoSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	keyword := r.URL.Query().Get("query")

	videos := services.Search(keyword)

	json.NewEncoder(w).Encode(videos)
}

func VideoDetail(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	videoURL, cookies, err := services.ViewVideoDetail(url)
	if err != nil {
		http.Error(w, "Failed to load video content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Accept-Ranges", "bytes")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", videoURL, nil)
	req.Header.Set("Cookie", cookies)

	rangeHeader := r.Header.Get("Range")

	var start, end int64
	fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", start))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to load content", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		http.Error(w, "Requested range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	contentLength := resp.ContentLength
	if end == 0 {
		end = contentLength - 1
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, contentLength))

	io.CopyN(w, resp.Body, end-start+1)
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-tiktok-scraping/models"
	"go-tiktok-scraping/services"
	"go-tiktok-scraping/utils"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	keyword := r.URL.Query().Get("q")

	if keyword == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Query parameter 'q' is required"})
		return
	}

	videos := services.Search(keyword)

	json.NewEncoder(w).Encode(videos)
}

func VideoDetail(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Query parameter 'url' is required"})
		return
	}

	video, err := services.ViewVideoDetail(url)
	if err != nil {
		http.Error(w, "Failed to load video content", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Accept-Ranges", "bytes")

	client := &http.Client{}

	req, _ := http.NewRequest("GET", video.Src, nil)
	req.Header.Set("Cookie", video.Cookie)

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

func Stream(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Query parameter 'q' is required"})
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Transfer-Encoding", "chunked")

	chromeManager := utils.GetChromeManager()
	allocCtx := chromeManager.GetContext()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if _, ok := ev.(*page.EventFrameSubtreeWillBeDetached); ok {
			log.Println("Ignoring frame detach event...")
		}
	})

	ctx, cancel = context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	url := "https://www.tiktok.com/search?q=" + keyword
	log.Println("Navigating to:", url)
	startUrl := time.Now()

	selector := `div[data-e2e="search_top-item"]`
	var links []string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		scrollAndLoadMoreContent(0),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('div[data-e2e="search_top-item"] a')).map(m=>m.href)`, &links),
	)
	if err != nil {
		log.Println("Failed to retrieve links:", err)
		return
	}
	log.Println("Total url:", len(links))
	elapsedUrl := time.Since(startUrl)
	fmt.Printf("Get all url time: %s\n", elapsedUrl)

	chunks := make(chan models.Video)
	linksChan := make(chan string)

	// Create a limited number of workers
	const concurrentWorkers = 3
	var wg sync.WaitGroup

	// Worker function
	worker := func() {
		defer wg.Done()
		for link := range linksChan {
			startDetail := time.Now()
			video, err := services.ViewVideoDetail(link)
			elapsedDetail := time.Since(startDetail)
			fmt.Printf("Get video detail time: %s\n", elapsedDetail)
			if err != nil {
				log.Println("Error retrieving video source for link:", link, "Error:", err)
				continue
			}
			chunks <- video
		}
	}

	// Start workers
	for i := 0; i < concurrentWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	// Send links to workers
	go func() {
		for _, link := range links {
			linksChan <- link
		}
		close(linksChan) // Close the channel after sending all URLs
	}()

	// Wait for workers to finish and close chunks
	go func() {
		wg.Wait()
		close(chunks)
	}()

	// Stream results
	for chunk := range chunks {
		json.NewEncoder(w).Encode(chunk)
		flusher.Flush()
	}
}

func scrollAndLoadMoreContent(scrolls int) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		for i := 0; i < scrolls; i++ {
			err := chromedp.Run(ctx, chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`, nil))
			if err != nil {
				return err
			}
			time.Sleep(2 * time.Second)
		}
		return nil
	})
}

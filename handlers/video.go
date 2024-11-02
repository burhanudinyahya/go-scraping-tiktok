package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-tiktok-scraping/services"
	"go-tiktok-scraping/utils"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
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

// Handler function to send chunked data
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

	// Set headers for chunked response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Set up Chromedp options
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("user-agent", utils.GetRandomUserAgent()),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Ignore unnecessary events to avoid spamming logs
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if _, ok := ev.(*page.EventFrameSubtreeWillBeDetached); ok {
			log.Println("Ignoring frame detach event...")
		}
	})

	// Set a timeout for the entire context
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := "https://www.tiktok.com/search?q=" + keyword
	log.Println("Navigating to:", url)

	// Check if selector is available
	selector := `div[data-e2e="search_top-item"]`
	var links []string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('div[data-e2e="search_top-item"] a')).map(m=>m.href)`, &links),
	)
	if err != nil {
		log.Println("Failed to retrieve links:", err)
		return
	}

	// Prepare for chunked response with goroutines
	chunks := make(chan string)
	var wg sync.WaitGroup

	// Variables to hold scraped data
	var videoSrc, coverSrc string
	var cookies []*network.Cookie

	// Start goroutines to process each link
	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			// Run Chromedp commands and check for errors
			err := chromedp.Run(ctx,
				network.Enable(),
				chromedp.Navigate(link),
				chromedp.ActionFunc(func(ctx context.Context) error {
					var err error
					cookies, err = network.GetCookies().Do(ctx)
					return err
				}),
				chromedp.WaitVisible(`video`, chromedp.ByQuery),
				chromedp.AttributeValue(`video > source`, "src", &videoSrc, nil),
			)
			if err != nil {
				log.Println("Error retrieving video source for link:", link, "Error:", err)
				return // Skip sending this chunk due to error
			}

			// Convert cookies to a single string format
			cookieStrings := make([]string, len(cookies))
			for i, cookie := range cookies {
				cookieStrings[i] = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
			}
			cookiesJoined := strings.Join(cookieStrings, "; ")

			// Format chunk data only if chromedp.Run succeeded
			chunk := fmt.Sprintf(`{
				"id": "%s",
				"title": "%s",
				"tags": [%s],
				"url": "%s",
				"cover": "%s",
				"src": "%s",
				"username": "%s",
				"cookies": "%s"
			}`,
				utils.GetID(link),
				"", // Title placeholder
				strings.Join(utils.ExtractTags(link), ", "),
				videoSrc,
				coverSrc,
				videoSrc,
				utils.GetUsername(link),
				cookiesJoined,
			)

			// Send chunk to channel only if no error occurred
			chunks <- chunk
		}(link)
	}

	// Close channel after all goroutines complete
	go func() {
		wg.Wait()
		close(chunks)
	}()

	// Stream each chunk to the client as it arrives
	for chunk := range chunks {
		fmt.Fprintf(w, "%s\n", chunk)
		flusher.Flush() // Flush to send chunk to client immediately
	}
}

package services

import (
	"context"
	"go-tiktok-scraping/models"
	"log"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func Search(keyword string) []models.Video {

	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if _, ok := ev.(*page.EventFrameSubtreeWillBeDetached); ok {
			log.Println("Ignoring frame detach event...")
		}
	})

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := "https://www.tiktok.com/search?q=" + keyword
	var videoLinks []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div[data-e2e="search_top-item"]`, chromedp.ByQuery),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('div[data-e2e="search_top-item"] a')).map(a => a.href)`, &videoLinks),
	)
	if err != nil {
		log.Fatal(err)
	}

	var videos []models.Video
	for i, link := range videoLinks {
		videos = append(videos, models.Video{ID: i, URL: "http://localhost:8000/video?url=" + link})
	}
	return videos
}

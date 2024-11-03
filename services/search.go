package services

import (
	"context"
	"fmt"
	"go-tiktok-scraping/models"
	"go-tiktok-scraping/utils"
	"log"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func fetchVideos(keyword string) ([]models.Video, error) {

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
	log.Println(url)

	selector := `div[data-e2e="search_top-item"]`

	// cek dulu selector ada ga
	var links []string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('div[data-e2e="search_top-item"] a')).map(m=>m.href)`, &links),
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(links)

	var videoLinks []*models.Video

	// fmt.Println(selector)
	querySeactor := `Array.from(document.querySelectorAll('div[data-e2e="search_top-item-list"] > div')).map(item => ({
		url: item.querySelector('div[data-e2e="search_top-item"] a').href,
		src: item.querySelector('div[data-e2e="search_top-item"] picture > img').src,
		title: item.querySelector('div[data-e2e="search-card-desc"] h1').innerText,
	}))`
	// fmt.Println(querySeactor)
	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Evaluate(querySeactor, &videoLinks),
		// chromedp.Evaluate(`Array.from(document.querySelectorAll('`+selector+` a')).map(a => a.href)`, &videoLinks),
	)
	if err != nil {
		return nil, err
	}

	// videos := make([]models.Video, len(videoLinks))
	// for i, link := range videoLinks {
	// 	videos[i] = models.Video{ID: i, URL: "http://localhost:8000/video?url=" + link}
	// }

	var videos []models.Video
	for _, videoPtr := range videoLinks {
		if videoPtr != nil { // Check for nil pointer
			video := models.Video{
				ID:       utils.GetID(videoPtr.URL),
				Title:    videoPtr.Title,
				Tags:     utils.ExtractTags(videoPtr.Title),
				URL:      videoPtr.URL,
				Src:      videoPtr.Src,
				Username: utils.GetUsername(videoPtr.URL),
			}
			videos = append(videos, video) // Append to the slice
		}
	}
	return videos, nil
}

func Search(keyword string) []models.Video {
	videos, err := fetchVideos(keyword)
	if err != nil {
		log.Println(err)
	}
	return videos
}

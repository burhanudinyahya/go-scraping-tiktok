package services

import (
	"context"
	"fmt"
	"go-tiktok-scraping/models"
	"go-tiktok-scraping/utils"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func ViewVideoDetail(url string) (models.Video, error) {

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

	var videoSrc, image, title, likes, comments string
	var cookies []*network.Cookie

	err := chromedp.Run(ctx,
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
		chromedp.WaitVisible(`video`, chromedp.ByQuery),
		chromedp.AttributeValue(`video > source`, "src", &videoSrc, nil),
		chromedp.WaitVisible(`picture`, chromedp.ByQuery),
		chromedp.AttributeValue(`picture > img`, "src", &image, nil),
		chromedp.WaitVisible(`h1`, chromedp.ByQuery),
		chromedp.Evaluate(`document.querySelector("h1").innerText`, &title),
		chromedp.WaitVisible(`strong[data-e2e="like-count"]`),
		chromedp.Text(`strong[data-e2e="like-count"]`, &likes),
		chromedp.WaitVisible(`strong[data-e2e="comment-count"]`),
		chromedp.Text(`strong[data-e2e="comment-count"]`, &comments),
	)
	if err != nil {
		log.Println(err)
	}

	cookieStrings := make([]string, len(cookies))
	for index, cookie := range cookies {
		cookieStrings[index] = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
	}

	cookiesJoined := strings.Join(cookieStrings, "; ")

	video := models.Video{
		ID:       utils.GetID(url),
		Title:    title,
		Tags:     utils.ExtractTags(title),
		URL:      url,
		Src:      videoSrc,
		Image:    image,
		Username: utils.GetUsername(url),
		Cookie:   cookiesJoined,
		Likes:    likes,
		Comments: comments,
	}

	return video, nil
}

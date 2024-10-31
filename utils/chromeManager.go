package utils

import (
	"context"
	"fmt"
	"sync"

	"github.com/chromedp/chromedp"
)

type ChromeManager struct {
	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc
}

func (c *ChromeManager) Initialize() {
	c.once.Do(func() {
		opts := []chromedp.ExecAllocatorOption{
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("ignore-certificate-errors", true),
		}

		allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel := chromedp.NewContext(allocCtx)
		c.ctx = ctx
		c.cancel = cancel
		fmt.Println("Chromedp session started")
	})
}

// GetContext returns the existing Chrome context
func (c *ChromeManager) GetContext() context.Context {
	c.Initialize()
	return c.ctx
}

func (c *ChromeManager) Close() {
	if c.cancel != nil {
		c.cancel()
	}
}

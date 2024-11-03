package utils

import (
	"context"
	"fmt"
	"sync"

	"github.com/chromedp/chromedp"
)

// ChromeManager is a struct for managing a Chrome session
type ChromeManager struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var instance *ChromeManager
var once sync.Once

// GetChromeManager returns the singleton instance of ChromeManager
func GetChromeManager() *ChromeManager {
	once.Do(func() {
		instance = &ChromeManager{}
		instance.Initialize()
	})
	return instance
}

func (c *ChromeManager) Initialize() {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("user-agent", GetRandomUserAgent()),
		// chromedp.Flag("proxy-server", GetRandomProxyServers()),
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	c.ctx = ctx
	c.cancel = cancel
	fmt.Println("Chromedp session started")
}

// GetContext returns the existing Chrome context
func (c *ChromeManager) GetContext() context.Context {
	return c.ctx
}

// Close closes the Chrome session
func (c *ChromeManager) Close() {
	if c.cancel != nil {
		c.cancel()
	}
}

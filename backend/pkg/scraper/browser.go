package scraper

import (
	"context"

	"github.com/chromedp/chromedp"
)

// BrowserOptions holds configuration for the headless browser.
type BrowserOptions struct {
	// Headless defines whether the browser should run without a GUI.
	// Set to false for local debugging to see the browser in action.
	Headless bool

	// UserAgent to use for HTTP requests to avoid bot detection.
	UserAgent string
}

// DefaultBrowserOptions returns the recommended options for scraping.
func DefaultBrowserOptions() BrowserOptions {
	return BrowserOptions{
		Headless:  true,
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	}
}

// NewBrowserContext creates a new chromedp context (a browser tab) with the given options.
// The caller is responsible for calling both cancel functions:
//
//	allocCtx, allocCancel := scraper.NewBrowserContext(ctx, opts)
//	defer allocCancel()
func NewBrowserContext(ctx context.Context, opts BrowserOptions) (context.Context, context.CancelFunc) {
	allocOpts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", opts.Headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent(opts.UserAgent),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, allocOpts...)
	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)

	// Combine cancel functions so the caller only needs one
	return browserCtx, func() {
		cancelBrowser()
		cancel()
	}
}

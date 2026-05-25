package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/model"
)

type AlfagiftScraper struct{}

func NewAlfagiftScraper() *AlfagiftScraper {
	return &AlfagiftScraper{}
}

func (s *AlfagiftScraper) ExtractProducts(ctx context.Context, targetURL string) ([]model.Product, error) {
	log.Printf("[AlfagiftScraper] Starting browser...")
	
	// Prepare options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"),
	)
	
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	taskCtx, cancelTask := chromedp.NewContext(allocCtx)
	defer cancelTask()

	taskCtx, cancelTimeout := context.WithTimeout(taskCtx, 60*time.Second)
	defer cancelTimeout()

	log.Printf("[AlfagiftScraper] Navigating to %s...", targetURL)

	var products []model.Product
	
	// We use chromedp to evaluate a fetch request to the internal API, which bypasses the frontend
	// location modal restrictions while reusing the browser's dynamic tokens.
	fetchScript := `
		window.apiResponse = null;
		fetch('https://webcommerce-gw.alfagift.id/v2/products/category/5b85712ca3834cdebbbc4363?sortDirection=asc&start=0&limit=60', {
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json'
			}
		}).then(res => res.json()).then(data => { window.apiResponse = data; });
	`
	extractScript := `
		(() => {
			let results = [];
			if (window.apiResponse && window.apiResponse.products) {
				window.apiResponse.products.forEach(p => {
					results.push({
						name: p.productName,
						price: p.finalPrice ? p.finalPrice : p.basePrice,
						image_url: p.image
					});
				});
			}
			return results;
		})()
	`

	var rawResults []struct {
		Name     string  `json:"name"`
		Price    float64 `json:"price"`
		ImageURL string  `json:"image_url"`
	}

	err := chromedp.Run(taskCtx,
		chromedp.Navigate("https://alfagift.id/"),
		// Wait for the main app to load and set up fingerprint tokens
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(fetchScript, nil),
		chromedp.Poll(`window.apiResponse !== null`, nil),
		chromedp.Evaluate(extractScript, &rawResults),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract products from Alfagift API: %w", err)
	}

	log.Printf("[AlfagiftScraper] Found %d raw items via API.", len(rawResults))

	for _, item := range rawResults {
		// Even if price is 0 (due to location), we still save the product structure
		products = append(products, model.Product{
			Name:        item.Name,
			Price:       item.Price,
			SourceURL:   targetURL,
			ImageURL:    item.ImageURL,
			Supermarket: "Alfagift",
		})
	}

	return products, nil
}

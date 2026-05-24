package scraper

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
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
	
	// We use a robust JS evaluation to extract products since CSS classes might change.
	// It scans common e-commerce product card structures.
	extractScript := `
		(() => {
			let results = [];
			// Try to find product cards based on common keywords in classes
			let cards = document.querySelectorAll('[class*="product-card"], [class*="productCard"], .list-product-catalog > div, [class*="product-item"]');
			
			cards.forEach(card => {
				let nameEl = card.querySelector('[class*="name"], [class*="title"], p.fw7, p.fw5');
				let priceEl = card.querySelector('[class*="price"], [class*="text-danger"], p.text-primary, p[class*="price"]');
				let imgEl = card.querySelector('img');
				
				if (nameEl && priceEl && imgEl) {
					let name = nameEl.innerText.trim();
					let priceStr = priceEl.innerText.trim();
					let image = imgEl.src;
					
					// Basic validation
					if (name.length > 0 && priceStr.includes('Rp')) {
						results.push({
							name: name,
							price: priceStr,
							image_url: image
						});
					}
				}
			});
			return results;
		})()
	`

	var rawResults []struct {
		Name     string `json:"name"`
		Price    string `json:"price"`
		ImageURL string `json:"image_url"`
	}

	err := chromedp.Run(taskCtx,
		chromedp.Navigate(targetURL),
		// Wait for the main app to load
		chromedp.Sleep(10*time.Second),
		
		// Attempt to click the location modal overlay or default button if it exists
		// In Alfagift, it often has a default "JABODETABEK" overlay. We try clicking any button that says "Pilih" or "Masuk"
		chromedp.ActionFunc(func(c context.Context) error {
			// This is a best-effort click to dismiss modals
			ctx, cancel := context.WithTimeout(c, 2*time.Second)
			defer cancel()
			var nodes []*cdp.Node
			_ = chromedp.Nodes(`button`, &nodes, chromedp.ByQueryAll).Do(ctx)
			for _, n := range nodes {
				var text string
				chromedp.Text([]cdp.NodeID{n.NodeID}, &text, chromedp.ByNodeID).Do(ctx)
				if strings.Contains(strings.ToLower(text), "pilih") || strings.Contains(strings.ToLower(text), "mengerti") {
					chromedp.MouseClickNode(n).Do(ctx)
					time.Sleep(1 * time.Second)
				}
			}
			return nil
		}),
		
		// Wait again for products to fetch
		chromedp.Sleep(5*time.Second),
		
		// Extract using JS
		chromedp.Evaluate(extractScript, &rawResults),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract products from Alfagift: %w", err)
	}

	log.Printf("[AlfagiftScraper] Found %d raw items.", len(rawResults))

	for _, item := range rawResults {
		// Clean up price: "Rp 15.500" -> 15500
		cleanedPrice := strings.ReplaceAll(item.Price, "Rp", "")
		cleanedPrice = strings.ReplaceAll(cleanedPrice, ".", "")
		cleanedPrice = strings.ReplaceAll(cleanedPrice, " ", "")
		cleanedPrice = strings.TrimSpace(cleanedPrice)

		priceInt, err := strconv.ParseFloat(cleanedPrice, 64)
		if err != nil {
			log.Printf("[AlfagiftScraper] Warning: failed to parse price '%s' for item '%s'", item.Price, item.Name)
			continue
		}

		products = append(products, model.Product{
			Name:        item.Name,
			Price:       priceInt,
			SourceURL:   targetURL,
			ImageURL:    item.ImageURL,
			Supermarket: "Alfagift",
		})
	}

	return products, nil
}

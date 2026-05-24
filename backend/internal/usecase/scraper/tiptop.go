package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/model"
	pkgscraper "github.com/mhamdriizki/grocery-scrapping-automation/backend/pkg/scraper"
)

const (
	tipTopBaseURL    = "https://shop.tiptop.co.id/outlet/Ciputat"
	tipTopSupermarket = "Tip Top"
	tipTopCategory   = "Keperluan Dapur"
)

// TipTopScraper handles the scraping logic for the Tip Top Supermarket website.
type TipTopScraper struct{}

// NewTipTopScraper creates a new instance of TipTopScraper.
func NewTipTopScraper() *TipTopScraper {
	return &TipTopScraper{}
}

// ScrapedProduct is a raw intermediate struct used during HTML parsing.
type ScrapedProduct struct {
	Name     string
	Price    string
	ImageURL string
}

// ScrapeKepDateperDapur navigates to the Tip Top Ciputat outlet page,
// clicks into the "Keperluan Dapur" category, and extracts product data.
// It prints each product to the terminal (Phase 1 - no DB storage).
func (s *TipTopScraper) ScrapeKepDateperDapur(ctx context.Context) ([]model.Product, error) {
	log.Println("[TipTopScraper] Starting browser...")

	opts := pkgscraper.DefaultBrowserOptions()
	browserCtx, cancel := pkgscraper.NewBrowserContext(ctx, opts)
	defer cancel()

	var rawProducts []ScrapedProduct

	// Build chromedp action sequence
	err := chromedp.Run(browserCtx,
		// Step 1: Navigate to outlet page
		chromedp.Navigate(tipTopBaseURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Printf("[TipTopScraper] Navigating to %s... done", tipTopBaseURL)
			return nil
		}),

		// Step 2: Wait for category list to be visible (SPA renders async)
		chromedp.WaitVisible(`a.category-item, .category-link, a[href*="category"]`, chromedp.ByQuery),

		// Step 3: Find and click the "Keperluan Dapur" category link
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[TipTopScraper] Looking for 'Keperluan Dapur' category...")
			return chromedp.Click(`//a[contains(translate(., 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', 'abcdefghijklmnopqrstuvwxyz'), 'keperluan dapur')]`, chromedp.BySearch).Do(ctx)
		}),

		// Step 4: Wait for product list to load after category click
		// The SPA might show skeleton loaders with text "...Loading" initially.
		// We use Poll to wait until the product text is no longer "Loading".
		chromedp.Poll(`
			(function() {
				const cards = document.querySelectorAll('.product-item, .product-card, [class*="product"]');
				if (cards.length === 0) return false;
				for (let i = 0; i < cards.length; i++) {
					const text = cards[i].innerText;
					if (text && text.includes('Loading')) {
						return false; // Still loading
					}
				}
				// Also ensure we actually have some text
				return cards[0].innerText.trim().length > 0;
			})()
		`, nil, chromedp.WithPollingInterval(1*time.Second), chromedp.WithPollingTimeout(30*time.Second)),
		// Sleep a bit more just to be safe that all elements are fully populated
		chromedp.Sleep(2*time.Second),

		// Step 5: Extract all product data via JavaScript
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("[TipTopScraper] Extracting products via JavaScript...")

			// Use JS to extract product name, price, and image from the DOM
			var result []map[string]interface{}
			err := chromedp.Evaluate(`
				(function() {
					const products = [];
					// Try multiple common selectors for product cards
					const cards = document.querySelectorAll('.product-item, .product-card, [class*="product"]');
					cards.forEach(card => {
						const nameEl = card.querySelector('[class*="name"], [class*="title"], h3, h4, h5');
						const priceEl = card.querySelector('[class*="price"], [class*="harga"]');
						const imgEl = card.querySelector('img');
						if (nameEl && priceEl) {
							products.push({
								name: nameEl.innerText.trim(),
								price: priceEl.innerText.trim(),
								image: imgEl ? imgEl.src : '',
							});
						}
					});
					return products;
				})()
			`, &result).Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to evaluate JS for products: %w", err)
			}

			for _, item := range result {
				name, _ := item["name"].(string)
				price, _ := item["price"].(string)
				image, _ := item["image"].(string)
				if name != "" {
					rawProducts = append(rawProducts, ScrapedProduct{
						Name:     name,
						Price:    price,
						ImageURL: image,
					})
				}
			}
			return nil
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("chromedp run failed: %w", err)
	}

	// Convert raw scraped data to model.Product and print to terminal
	scrapedAt := time.Now()
	var products []model.Product

	log.Printf("[TipTopScraper] Found %d products in '%s' category:", len(rawProducts), tipTopCategory)
	log.Println(strings.Repeat("-", 60))

	for i, raw := range rawProducts {
		price := parsePrice(raw.Price)
		p := model.Product{
			Name:        raw.Name,
			Price:       price,
			ImageURL:    raw.ImageURL,
			SourceURL:   tipTopBaseURL,
			Category:    tipTopCategory,
			Supermarket: tipTopSupermarket,
			ScrapedAt:   scrapedAt,
		}
		products = append(products, p)

		// Print each product to terminal
		log.Printf("[%d] Name: %-50s | Price: Rp %.0f", i+1, raw.Name, price)
	}

	log.Println(strings.Repeat("-", 60))
	log.Printf("[TipTopScraper] Scraping completed. Total: %d products.", len(products))

	return products, nil
}

// parsePrice converts a raw price string like "Rp 12.500" or "12500" to float64.
func parsePrice(raw string) float64 {
	// Remove common non-numeric characters
	cleaned := strings.ReplaceAll(raw, "Rp", "")
	cleaned = strings.ReplaceAll(cleaned, "rp", "")
	cleaned = strings.ReplaceAll(cleaned, ".", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)

	var price float64
	fmt.Sscanf(cleaned, "%f", &price)
	return price
}

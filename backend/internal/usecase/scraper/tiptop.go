package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/model"
)

const (
	tipTopSupermarket = "Tip Top"
	tipTopCategory    = "Keperluan Dapur"
	tipTopOutletID    = "63c4b0eceb09e61d3a2eeb4d" // Ciputat
	tipTopCategoryID  = "63b9444d9121c343a7d3cbc7" // Keperluan Dapur
	tipTopAPIKey      = "29ea3dcf-67c3-45f6-b5c4-d2628f7e09fa"
)

// TipTopScraper handles the scraping logic for the Tip Top Supermarket website.
// Instead of a flaky headless browser, it leverages TipTop's internal REST API for 100% reliability.
type TipTopScraper struct{}

// NewTipTopScraper creates a new instance of TipTopScraper.
func NewTipTopScraper() *TipTopScraper {
	return &TipTopScraper{}
}

// ScrapeKeperluanDapur fetches product data for Ciputat / Keperluan Dapur via internal API.
func (s *TipTopScraper) ScrapeKeperluanDapur(ctx context.Context) ([]model.Product, error) {
	log.Println("[TipTopScraper] Starting API extraction...")

	// Create request
	url := fmt.Sprintf("https://api.tiptop.co.id/api/web/product?limit=100&categoryId=%s&page=1&outletId=%s", tipTopCategoryID, tipTopOutletID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add required headers discovered during network analysis
	req.Header.Set("x-api-key", tipTopAPIKey)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON
	var result struct {
		Status     bool   `json:"status"`
		StatusCode int    `json:"statusCode"`
		Message    string `json:"message"`
		Data       []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			MainImage   struct {
				SmallImage string `json:"small_image"`
				LargeImage string `json:"large_image"`
			} `json:"main_image"`
			InfoProduct struct {
				PricingStock []struct {
					Name               string `json:"name"`
					PricingStockOutlet []struct {
						Price        float64 `json:"price"`
						SpecialPrice float64 `json:"special_price"`
					} `json:"pricing_stock_outlet"`
				} `json:"pricing_stock"`
			} `json:"info_product"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var products []model.Product
	scrapedAt := time.Now()

	log.Printf("[TipTopScraper] Found %d products from API", len(result.Data))
	log.Println(strings.Repeat("-", 80))

	for i, item := range result.Data {
		var finalPrice float64
		// Extract price from the first variant
		if len(item.InfoProduct.PricingStock) > 0 && len(item.InfoProduct.PricingStock[0].PricingStockOutlet) > 0 {
			stock := item.InfoProduct.PricingStock[0].PricingStockOutlet[0]
			finalPrice = stock.Price
			if stock.SpecialPrice > 0 && stock.SpecialPrice < stock.Price {
				finalPrice = stock.SpecialPrice
			}
		}

		p := model.Product{
			Name:        item.Name,
			Price:       finalPrice,
			ImageURL:    item.MainImage.SmallImage,
			SourceURL:   url, // keeping the API URL as source for traceability
			Category:    tipTopCategory,
			Supermarket: tipTopSupermarket,
			ScrapedAt:   scrapedAt,
		}
		products = append(products, p)

		log.Printf("[%d] Name: %-50s | Price: Rp %.0f", i+1, p.Name, p.Price)
	}

	log.Println(strings.Repeat("-", 80))
	log.Printf("[TipTopScraper] API Extraction completed. Total: %d products.", len(products))

	return products, nil
}

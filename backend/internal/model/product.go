package model

import "time"

// Product represents a grocery product scraped from a supermarket website.
type Product struct {
	BaseModel

	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
	SourceURL  string  `json:"source_url"`
	Category   string  `json:"category"`
	Supermarket string `json:"supermarket"`
	ScrapedAt  time.Time `json:"scraped_at"`
}

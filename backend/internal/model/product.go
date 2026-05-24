package model

import "time"

// Product represents a grocery product scraped from a supermarket website.
type Product struct {
	BaseModel

	Name        string    `json:"name" gorm:"uniqueIndex:idx_name_supermarket"`
	Price       float64   `json:"price"`
	ImageURL    string    `json:"image_url"`
	SourceURL   string    `json:"source_url"`
	Category    string    `json:"category"`
	Supermarket string    `json:"supermarket" gorm:"uniqueIndex:idx_name_supermarket"`
	ScrapedAt   time.Time `json:"scraped_at"`
}

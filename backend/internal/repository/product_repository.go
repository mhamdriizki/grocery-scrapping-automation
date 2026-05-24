package repository

import (
	"context"
	"fmt"

	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductRepository defines the contract for product database operations.
type ProductRepository interface {
	Save(ctx context.Context, product *model.Product) error
	SaveBatch(ctx context.Context, products []model.Product) error
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new instance of ProductRepository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// Save inserts a single product or updates its details if it already exists (upsert).
func (r *productRepository) Save(ctx context.Context, product *model.Product) error {
	// Upsert logic: On conflict (Name + Supermarket), update Price, ImageURL, SourceURL, ScrapedAt.
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}, {Name: "supermarket"}}, // unique index
		DoUpdates: clause.AssignmentColumns([]string{
			"price", "image_url", "source_url", "scraped_at", "updated_at",
		}),
	}).Create(product).Error

	if err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}
	return nil
}

// SaveBatch inserts multiple products or updates them if they already exist (upsert) in batches.
func (r *productRepository) SaveBatch(ctx context.Context, products []model.Product) error {
	if len(products) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}, {Name: "supermarket"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"price", "image_url", "source_url", "scraped_at", "updated_at",
		}),
	}).CreateInBatches(&products, 100).Error

	if err != nil {
		return fmt.Errorf("failed to save batch products: %w", err)
	}
	return nil
}

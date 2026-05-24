package product

import (
	"context"

	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/model"
	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/repository"
)

// ProductUsecase defines the interface for product-related business logic
type ProductUsecase interface {
	SearchProducts(ctx context.Context, searchQuery string) ([]model.Product, error)
}

type productUsecase struct {
	productRepo repository.ProductRepository
}

// NewProductUsecase creates a new ProductUsecase instance
func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return &productUsecase{
		productRepo: repo,
	}
}

// SearchProducts returns a list of products matching the search query
func (u *productUsecase) SearchProducts(ctx context.Context, searchQuery string) ([]model.Product, error) {
	return u.productRepo.FindAll(ctx, searchQuery)
}

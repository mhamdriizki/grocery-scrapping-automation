package http

import (
	"net/http"

	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/usecase/product"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	productUsecase product.ProductUsecase
}

// NewProductHandler creates a new instance of ProductHandler
func NewProductHandler(usecase product.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: usecase,
	}
}

// GetProducts handles GET /api/v1/products
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	search := r.URL.Query().Get("search")

	products, err := h.productUsecase.SearchProducts(r.Context(), search)
	if err != nil {
		JSON(w, http.StatusInternalServerError, false, "Failed to retrieve products", nil)
		return
	}

	JSON(w, http.StatusOK, true, "Success retrieve products", products)
}

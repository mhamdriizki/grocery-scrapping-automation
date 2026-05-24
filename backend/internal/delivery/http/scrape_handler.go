package http

import (
	"encoding/json"
	"net/http"

	"github.com/mhamdriizki/grocery-scrapping-automation/backend/internal/worker"
)

// ScrapeHandler handles HTTP requests to trigger scraper
type ScrapeHandler struct {
	distributor worker.TaskDistributor
}

// NewScrapeHandler creates a new instance of ScrapeHandler
func NewScrapeHandler(distributor worker.TaskDistributor) *ScrapeHandler {
	return &ScrapeHandler{
		distributor: distributor,
	}
}

// ScrapeRequest represents the expected payload for POST /api/v1/scrape
type ScrapeRequest struct {
	TargetURL string `json:"target_url"`
}

// TriggerScrape handles POST /api/v1/scrape
func (h *ScrapeHandler) TriggerScrape(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSON(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
		return
	}

	var req ScrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSON(w, http.StatusBadRequest, false, "Invalid JSON payload", nil)
		return
	}

	// For TipTop specifically as requested
	if req.TargetURL == "" {
		req.TargetURL = "https://shop.tiptop.co.id/outlet/Ciputat/category/Keperluan-Dapur?key=63b9444d9121c343a7d3cbc7&item=63c34ab03ac2ba06639c0b36"
	}

	err := h.distributor.DistributeScrapeGroceryTask(r.Context(), worker.ScrapeGroceryPayload{
		TargetURL: req.TargetURL,
	})

	if err != nil {
		JSON(w, http.StatusInternalServerError, false, "Failed to enqueue scrape task", nil)
		return
	}

	JSON(w, http.StatusOK, true, "Successfully enqueued scrape task", req)
}

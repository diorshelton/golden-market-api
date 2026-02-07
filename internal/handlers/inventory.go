package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/diorshelton/golden-market-api/internal/middleware"
	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
)

type InventoryServiceInterface interface {
	GetUserInventory(ctx context.Context, userID uuid.UUID) ([]models.InventoryItemDetail, error)
}

type InventoryHandler struct {
	inventoryService InventoryServiceInterface
}

func NewInventoryHandler(service InventoryServiceInterface) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: service,
	}
}

// GetInventory handles GET /api/v1/inventory
func (h *InventoryHandler) GetInventory(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	items, err := h.inventoryService.GetUserInventory(r.Context(), userID)
	if err != nil {
		log.Printf("GetInventory error for user %s: %v", userID, err)
		http.Error(w, "failed to get inventory", http.StatusInternalServerError)
		return
	}

	// Return empty array instead of null
	if items == nil {
		items = []models.InventoryItemDetail{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

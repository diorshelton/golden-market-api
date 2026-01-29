package inventory

import (
	"context"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/diorshelton/golden-market-api/internal/repository"
	"github.com/google/uuid"
)

type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
}

func NewInventoryService(inventoryRepo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
	}
}

// GetUserInventory retrieves all inventory items for a user with product details
func (s *InventoryService) GetUserInventory(ctx context.Context, userID uuid.UUID) ([]models.InventoryItemDetail, error) {
	return s.inventoryRepo.GetByUserID(ctx, userID)
}

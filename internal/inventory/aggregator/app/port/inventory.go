package port

import (
	"context"

	"github.com/LudSkywalker/inventory-system/internal/inventory/aggregator/app/dto"
)

type GlobalInventoryUseCase interface {
	GetGlobalInventory(ctx context.Context, itemID, storeID string) (*dto.GlobalInventoryDTO, error)
	GetAllInventories(ctx context.Context) ([]*dto.GlobalInventoryDTO, error)
	GetStoreInventory(ctx context.Context, storeID string) ([]*dto.GlobalInventoryDTO, error)
	GetItemInventory(ctx context.Context, itemID string) ([]*dto.GlobalInventoryDTO, error)
}

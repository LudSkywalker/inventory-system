package port

import (
	"context"

	"github.com/LudSkywalker/inventory-system/internal/inventory/query/app/dto"
)

type InventoryQueryUseCase interface {
	GetAllInventories(ctx context.Context) ([]*dto.InventoryDTO, error)
	GetStoreInventory(ctx context.Context, storeID string) ([]*dto.InventoryDTO, error)
	GetItemInventory(ctx context.Context, itemID string) ([]*dto.InventoryDTO, error)
	GetInventory(ctx context.Context, query dto.InventoryQuery) (*dto.InventoryDTO, error)
}

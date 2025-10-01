package input

import (
	"context"

	"github.com/LudSkywalker/inventory-system/operator-service/app/dto"
)

type GlobalInventoryUseCase interface {
	ProcessInventoryEvent(ctx context.Context, event dto.GlobalInventoryDTO) error
	GetGlobalInventory(ctx context.Context, itemID, storeID string) (*dto.GlobalInventoryDTO, error)
	GetAllInventories(ctx context.Context) ([]*dto.GlobalInventoryDTO, error)
}

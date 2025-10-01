package port

import (
	"context"

	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/app/dto"
)

// InventoryUseCase defines the application interface for inventory operations
type InventoryUseCase interface {
	UpdateStock(ctx context.Context, cmd dto.UpdateStockCommand) error
	GetStock(ctx context.Context, query dto.GetStockQuery) (*dto.InventoryDTO, error)
	DeleteStock(ctx context.Context, cmd dto.DeleteStockCommand) error
}

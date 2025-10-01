package port

import (
	"context"

	"github.com/LudSkywalker/inventory-system/store-service/app/dto"
)

type InventoryUseCase interface {
	UpdateStock(ctx context.Context, cmd dto.UpdateStockCommand) error
	GetStock(ctx context.Context, query dto.GetStockQuery) (*dto.InventoryDTO, error)
	DeleteStock(ctx context.Context, cmd dto.DeleteStockCommand) error
}

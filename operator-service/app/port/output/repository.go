package output

import (
	"context"

	"github.com/LudSkywalker/inventory-system/operator-service/domain/entity"
)

type GlobalInventoryRepository interface {
	Save(ctx context.Context, inventory *entity.GlobalInventory) error
	FindByItemAndStore(ctx context.Context, itemID, storeID string) (*entity.GlobalInventory, error)
	FindAll(ctx context.Context) ([]*entity.GlobalInventory, error)
	Delete(ctx context.Context, itemID, storeID string) error
}

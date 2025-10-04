package port

import (
	"context"

	"github.com/LudSkywalker/inventory-system/store-service/domain/entity"
	"github.com/LudSkywalker/inventory-system/store-service/domain/event"
)

type InventoryRepository interface {
	Save(ctx context.Context, inventory *entity.Inventory) error
	Find(ctx context.Context, itemID, storeID string) (*entity.Inventory, error)
	List(ctx context.Context) ([]*entity.Inventory, error)
	Delete(ctx context.Context, itemID, storeID string) error
}

type EventPublisher interface {
	PublishInventoryChange(ctx context.Context, event event.InventoryEvent) error
}

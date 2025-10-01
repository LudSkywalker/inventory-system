package domain

import (
	"context"

	"github.com/LudSkywalker/inventory-system/internal/inventory/core/entity"
	"github.com/LudSkywalker/inventory-system/internal/inventory/core/event"
)

// Repository defines the interface for inventory persistence
type Repository interface {
	Save(ctx context.Context, inventory *entity.Inventory) error
	FindByItemAndStore(ctx context.Context, itemID, storeID string) (*entity.Inventory, error)
	Delete(ctx context.Context, itemID, storeID string) error
}

// EventPublisher defines the interface for publishing inventory events
type EventPublisher interface {
	PublishInventoryChange(ctx context.Context, event *event.InventoryEvent) error
}

// Service defines the domain service for local inventory management
type Service interface {
	UpdateInventory(ctx context.Context, inventory *entity.Inventory) error
	GetInventory(ctx context.Context, itemID, storeID string) (*entity.Inventory, error)
	DeleteInventory(ctx context.Context, itemID, storeID string) error
}

package domain

import (
	"context"
	"time"

	"github.com/LudSkywalker/inventory-system/internal/inventory/core/event"
	"github.com/LudSkywalker/inventory-system/internal/inventory/core/valueobject"
)

// Repository defines the interface for inventory persistence
type Repository interface {
	Save(ctx context.Context, inventory *Inventory) error
	Find(ctx context.Context, itemID, storeID string) (*Inventory, error)
	Delete(ctx context.Context, itemID, storeID string) error
}

// EventPublisher defines the interface for publishing inventory events
type EventPublisher interface {
	PublishInventoryChange(ctx context.Context, event event.InventoryEvent) error
}

// Inventory represents the local inventory aggregate
type Inventory struct {
	ItemID    string
	StoreID   string
	Quantity  valueobject.Quantity
	UpdatedAt time.Time
	Version   int64
}

// NewInventory creates a new Inventory instance
func NewInventory(itemID, storeID string, quantity valueobject.Quantity) *Inventory {
	return &Inventory{
		ItemID:    itemID,
		StoreID:   storeID,
		Quantity:  quantity,
		UpdatedAt: time.Now(),
		Version:   1,
	}
}

// UpdateQuantity updates the inventory quantity and timestamp
func (i *Inventory) UpdateQuantity(q valueobject.Quantity) error {
	if !q.IsValid() {
		return ErrInvalidQuantity
	}
	i.Quantity = q
	i.UpdatedAt = time.Now()
	i.Version++
	return nil
}

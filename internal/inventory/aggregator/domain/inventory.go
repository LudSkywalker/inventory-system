package domain

import (
	"context"
	"time"
)

// GlobalInventory represents the aggregated inventory state
type GlobalInventory struct {
	ItemID    string
	StoreID   string
	Quantity  int
	UpdatedAt time.Time
	Version   int64
}

// Repository defines the interface for global inventory persistence
type Repository interface {
	Save(ctx context.Context, inventory *GlobalInventory) error
	FindByItemAndStore(ctx context.Context, itemID, storeID string) (*GlobalInventory, error)
	FindAll(ctx context.Context) ([]*GlobalInventory, error)
	FindByStore(ctx context.Context, storeID string) ([]*GlobalInventory, error)
	FindByItem(ctx context.Context, itemID string) ([]*GlobalInventory, error)
}

func NewGlobalInventory(itemID, storeID string, quantity int) *GlobalInventory {
	return &GlobalInventory{
		ItemID:    itemID,
		StoreID:   storeID,
		Quantity:  quantity,
		UpdatedAt: time.Now(),
		Version:   1,
	}
}

func (g *GlobalInventory) UpdateQuantity(quantity int) {
	g.Quantity = quantity
	g.UpdatedAt = time.Now()
	g.Version++
}

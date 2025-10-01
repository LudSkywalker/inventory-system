package entity

import (
	"time"
)

// GlobalInventory represents the global state of an inventory item
type GlobalInventory struct {
	ItemID    string
	StoreID   string
	Quantity  int
	UpdatedAt time.Time
	Version   int64
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

package entity

import (
	"time"

	"github.com/LudSkywalker/inventory-system/internal/inventory/core/valueobject"
)

// Inventory represents the core inventory domain entity
type Inventory struct {
	ItemID    string
	StoreID   string
	Quantity  valueobject.Quantity
	UpdatedAt time.Time
}

func NewInventory(itemID, storeID string, quantity valueobject.Quantity) (*Inventory, error) {
	if itemID == "" || storeID == "" {
		return nil, ErrInvalidInput
	}

	return &Inventory{
		ItemID:    itemID,
		StoreID:   storeID,
		Quantity:  quantity,
		UpdatedAt: time.Now(),
	}, nil
}

func (i *Inventory) UpdateQuantity(q valueobject.Quantity) error {
	if !q.IsValid() {
		return ErrInvalidQuantity
	}
	i.Quantity = q
	i.UpdatedAt = time.Now()
	return nil
}

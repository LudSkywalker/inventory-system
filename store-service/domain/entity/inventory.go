package entity

import (
	"errors"
	"time"

	"github.com/LudSkywalker/inventory-system/store-service/domain/valueobject"
)

var (
	ErrInvalidQuantity = errors.New("invalid quantity")
)

type Inventory struct {
	ItemID    string
	StoreID   string
	Quantity  valueobject.Quantity
	UpdatedAt time.Time
}

func NewInventory(itemID, storeID string, quantity valueobject.Quantity) (*Inventory, error) {
	if itemID == "" || storeID == "" {
		return nil, errors.New("item ID and store ID are required")
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

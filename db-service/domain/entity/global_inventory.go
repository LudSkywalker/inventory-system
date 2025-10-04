package entity

import (
	"time"
)

// GlobalInventory represents the global state of an inventory item
type GlobalInventory struct {
	ItemID    string    `json:"item_id"`
	ItemName  string    `json:"item_name"`
	StoreID   string    `json:"store_id"`
	Quantity  int       `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int64     `json:"version"`
}

func NewGlobalInventory(itemID, itemName, storeID string, quantity int) *GlobalInventory {
	return &GlobalInventory{
		ItemID:    itemID,
		ItemName:  itemName,
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

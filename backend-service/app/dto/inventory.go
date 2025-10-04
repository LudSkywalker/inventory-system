package dto

// InventoryDTO represents inventory information
type InventoryDTO struct {
	ItemID    string `json:"item_id" example:"item-123"`
	ItemName  string `json:"item_name" example:"Item Name"`
	StoreID   string `json:"store_id" example:"store-456"`
	Quantity  int    `json:"quantity" example:"100"`
	UpdatedAt string `json:"updated_at" example:"2023-10-03T19:20:30Z"`
}

// StoreInventory represents inventory for a specific store
type StoreInventory struct {
	StoreID  string `json:"store_id" example:"store-456"`
	Quantity int    `json:"quantity" example:"100"`
}

// GroupedInventoryDTO represents inventory grouped by item
type GroupedInventoryDTO struct {
	ItemID        string           `json:"item_id" example:"item-123"`
	ItemName      string           `json:"item_name" example:"Item Name"`
	TotalQuantity int              `json:"total_quantity" example:"150"`
	Stores        []StoreInventory `json:"stores"`
}

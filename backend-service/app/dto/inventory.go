package dto

// InventoryDTO represents inventory information
type InventoryDTO struct {
	ItemID    string `json:"item_id" example:"item-123"`
	StoreID   string `json:"store_id" example:"store-456"`
	Quantity  int    `json:"quantity" example:"100"`
	UpdatedAt string `json:"updated_at" example:"2023-10-03T19:20:30Z"`
}

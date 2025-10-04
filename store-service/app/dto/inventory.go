package dto

// UpdateStockCommand represents a request to update stock quantity
type UpdateStockCommand struct {
	ItemID   string `json:"item_id" example:"item-123" binding:"required"`
	StoreID  string `json:"store_id" example:"store-456" binding:"required"`
	Quantity int    `json:"quantity" example:"100" binding:"required"`
}

// GetStockQuery represents a query to get stock information
type GetStockQuery struct {
	ItemID  string `json:"item_id" example:"item-123"`
	StoreID string `json:"store_id" example:"store-456"`
}

// DeleteStockCommand represents a request to delete stock
type DeleteStockCommand struct {
	ItemID  string `json:"item_id" example:"item-123"`
	StoreID string `json:"store_id" example:"store-456"`
}

// InventoryDTO represents inventory information
type InventoryDTO struct {
	ItemID    string `json:"item_id" example:"item-123"`
	StoreID   string `json:"store_id" example:"store-456"`
	Quantity  int    `json:"quantity" example:"100"`
	UpdatedAt string `json:"updated_at" example:"2023-10-03T19:20:30Z"`
}

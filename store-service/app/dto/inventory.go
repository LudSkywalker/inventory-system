package dto

// UpdateStockCommand represents a request to update stock quantity
type UpdateStockCommand struct {
	ItemID   string `json:"item_id" example:"0192a4b0-6b2a-7c4b-8d5e-123456789abc" binding:"required"`
	ItemName string `json:"item_name" example:"Item Name"`
	StoreID  string `json:"store_id" example:"0192a4b0-6b2a-7c4b-8d5e-123456789abd" binding:"required"`
	Quantity int    `json:"quantity" example:"100" binding:"required"`
}

// GetStockQuery represents a query to get stock information
type GetStockQuery struct {
	ItemID  string `json:"item_id" example:"0192a4b0-6b2a-7c4b-8d5e-123456789abc"`
	StoreID string `json:"store_id" example:"0192a4b0-6b2a-7c4b-8d5e-123456789abd"`
}

// DeleteStockCommand represents a request to delete stock
type DeleteStockCommand struct {
	ItemID  string `json:"item_id" example:"0192a4b0-6b2a-7c4b-8d5e-123456789abc"`
	StoreID string `json:"store_id" example:"0192a4b0-6b2a-7c4b-8d5e-123456789abd"`
}

// InventoryDTO represents inventory information
type InventoryDTO struct {
	ItemID    string `json:"item_id" example:"0192a4b0-6b2a-7c4b-8d5e-123456789abc"`
	ItemName  string `json:"item_name" example:"Item Name"`
	StoreID   string `json:"store_id" example:"0192a4b0-6b2a-7c4b-8d5e-123456789abd"`
	Quantity  int    `json:"quantity" example:"100"`
	UpdatedAt string `json:"updated_at" example:"2023-10-03T19:20:30Z"`
}

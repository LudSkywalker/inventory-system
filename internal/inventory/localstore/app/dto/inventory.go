package dto

type UpdateStockCommand struct {
	ItemID   string `json:"item_id"`
	StoreID  string `json:"store_id"`
	Quantity int    `json:"quantity"`
}

type GetStockQuery struct {
	ItemID  string `json:"item_id"`
	StoreID string `json:"store_id"`
}

type DeleteStockCommand struct {
	ItemID  string `json:"item_id"`
	StoreID string `json:"store_id"`
}

type InventoryDTO struct {
	ItemID    string `json:"item_id"`
	StoreID   string `json:"store_id"`
	Quantity  int    `json:"quantity"`
	UpdatedAt string `json:"updated_at"`
	Version   int64  `json:"version"`
}

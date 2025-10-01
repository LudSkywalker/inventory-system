package dto

type InventoryDTO struct {
	ItemID    string `json:"item_id"`
	StoreID   string `json:"store_id"`
	Quantity  int    `json:"quantity"`
	UpdatedAt string `json:"updated_at"`
	Version   int64  `json:"version"`
}

type InventoryQuery struct {
	ItemID  string `json:"item_id,omitempty"`
	StoreID string `json:"store_id,omitempty"`
}

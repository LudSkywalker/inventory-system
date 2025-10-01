package dto

type InventoryDTO struct {
	ItemID    string `json:"item_id"`
	StoreID   string `json:"store_id"`
	Quantity  int    `json:"quantity"`
	UpdatedAt string `json:"updated_at"`
}

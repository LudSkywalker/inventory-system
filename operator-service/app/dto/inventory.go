package dto

type GlobalInventoryDTO struct {
	ItemID    string `json:"item_id"`
	StoreID   string `json:"store_id"`
	Quantity  int    `json:"quantity"`
	UpdatedAt string `json:"updated_at"`
	Version   int64  `json:"version"`
}

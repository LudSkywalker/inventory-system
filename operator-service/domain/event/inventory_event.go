package event

import "time"

type InventoryEvent struct {
	EventID   string    `json:"event_id"`
	ItemID    string    `json:"item_id"`
	ItemName  string    `json:"item_name"`
	StoreID   string    `json:"store_id"`
	Quantity  int       `json:"quantity"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}

const (
	OperationUpdate = "UPDATE"
	OperationDelete = "DELETE"
)

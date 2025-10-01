package event

import "time"

// InventoryEvent represents any event related to inventory changes
type InventoryEvent struct {
	EventID   string    `json:"event_id"`
	ItemID    string    `json:"item_id"`
	StoreID   string    `json:"store_id"`
	Quantity  int       `json:"quantity"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
	Version   int64     `json:"version"`
}

// Operation types
const (
	OperationUpdate = "UPDATE"
	OperationDelete = "DELETE"
)

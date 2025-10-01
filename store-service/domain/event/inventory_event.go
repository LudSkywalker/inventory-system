package event

import "time"

type InventoryEvent struct {
	EventID   string    `json:"event_id"`
	ItemID    string    `json:"item_id"`
	StoreID   string    `json:"store_id"`
	Quantity  int       `json:"quantity"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}

func NewInventoryEvent(itemID, storeID string, quantity int, operation string) InventoryEvent {
	return InventoryEvent{
		EventID:   GenerateEventID(),
		ItemID:    itemID,
		StoreID:   storeID,
		Quantity:  quantity,
		Operation: operation,
		Timestamp: time.Now(),
	}
}

// Event types
const (
	OperationUpdate = "UPDATE"
	OperationDelete = "DELETE"
)

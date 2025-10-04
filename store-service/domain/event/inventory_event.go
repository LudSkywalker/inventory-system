package event

import (
	"time"

	"github.com/google/uuid"
)

type InventoryEvent struct {
	EventID   string    `json:"event_id"`
	ItemID    string    `json:"item_id"`
	ItemName  string    `json:"item_name"`
	StoreID   string    `json:"store_id"`
	Quantity  int       `json:"quantity"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
}

func NewInventoryEvent(itemID, itemName, storeID string, quantity int, operation string) InventoryEvent {
	return InventoryEvent{
		EventID:   GenerateEventID(),
		ItemID:    itemID,
		ItemName:  itemName,
		StoreID:   storeID,
		Quantity:  quantity,
		Operation: operation,
		Timestamp: time.Now(),
	}
}

func GenerateEventID() string {
	return uuid.Must(uuid.NewV7()).String()
}

// Event types
const (
	OperationUpdate = "UPDATE"
	OperationDelete = "DELETE"
)

package entity_test

import (
	"testing"
	"time"

	"github.com/LudSkywalker/inventory-system/store-service/domain/entity"
	"github.com/LudSkywalker/inventory-system/store-service/domain/valueobject"
)

func TestNewInventory(t *testing.T) {
	// Given
	itemID := "item1"
	itemName := "Item 1"
	storeID := "store1"
	quantity, _ := valueobject.NewQuantity(10)

	// When
	inventory, err := entity.NewInventory(itemID, itemName, storeID, quantity)

	// Then
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if inventory.ItemID != itemID {
		t.Errorf("Expected ItemID %s, got %s", itemID, inventory.ItemID)
	}

	if inventory.ItemName != itemName {
		t.Errorf("Expected ItemName %s, got %s", itemName, inventory.ItemName)
	}

	if inventory.StoreID != storeID {
		t.Errorf("Expected StoreID %s, got %s", storeID, inventory.StoreID)
	}

	if inventory.Quantity.Value() != 10 {
		t.Errorf("Expected Quantity 10, got %d", inventory.Quantity.Value())
	}

	if inventory.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestUpdateQuantity(t *testing.T) {
	// Given
	inventory, _ := entity.NewInventory("item1", "Item 1", "store1", valueobject.Quantity{})
	newQuantity, _ := valueobject.NewQuantity(20)
	oldUpdatedAt := inventory.UpdatedAt

	time.Sleep(1 * time.Millisecond) // Ensure time difference

	// When
	err := inventory.UpdateQuantity(newQuantity)

	// Then
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if inventory.Quantity.Value() != 20 {
		t.Errorf("Expected Quantity 20, got %d", inventory.Quantity.Value())
	}

	if !inventory.UpdatedAt.After(oldUpdatedAt) {
		t.Error("Expected UpdatedAt to be updated")
	}
}

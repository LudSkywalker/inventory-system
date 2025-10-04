package service_test

import (
	"context"
	"testing"

	"github.com/LudSkywalker/inventory-system/backend-service/app/dto"
	"github.com/LudSkywalker/inventory-system/backend-service/app/service"
)

type mockHTTPRepository struct {
	inventories []*dto.InventoryDTO
}

func (m *mockHTTPRepository) FindAll(ctx context.Context) ([]*dto.InventoryDTO, error) {
	return m.inventories, nil
}

func TestInventoryService_GetGroupedInventory(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockHTTPRepository{
		inventories: []*dto.InventoryDTO{
			{ItemID: "item1", ItemName: "Item 1", StoreID: "store1", Quantity: 10, UpdatedAt: "2023-10-01T10:00:00Z"},
			{ItemID: "item1", ItemName: "Item 1", StoreID: "store2", Quantity: 20, UpdatedAt: "2023-10-01T10:00:00Z"},
			{ItemID: "item2", ItemName: "Item 2", StoreID: "store1", Quantity: 5, UpdatedAt: "2023-10-01T10:00:00Z"},
		},
	}
	svc := &service.InventoryService{Repo: mockRepo}

	grouped, err := svc.GetGroupedInventory(ctx)
	if err != nil {
		t.Fatalf("GetGroupedInventory failed: %v", err)
	}

	if len(grouped) != 2 {
		t.Errorf("Expected 2 grouped items, got %d", len(grouped))
	}

	// Find items by ItemID since map iteration order is not guaranteed
	var item1, item2 dto.GroupedInventoryDTO
	for _, g := range grouped {
		if g.ItemID == "item1" {
			item1 = g
		} else if g.ItemID == "item2" {
			item2 = g
		}
	}

	// Check item1
	if item1.ItemID != "item1" || item1.ItemName != "Item 1" {
		t.Errorf("Item1 mismatch: %+v", item1)
	}
	if item1.TotalQuantity != 30 {
		t.Errorf("Expected total quantity 30, got %d", item1.TotalQuantity)
	}
	if len(item1.Stores) != 2 {
		t.Errorf("Expected 2 stores for item1, got %d", len(item1.Stores))
	}

	// Check item2
	if item2.ItemID != "item2" || item2.ItemName != "Item 2" {
		t.Errorf("Item2 mismatch: %+v", item2)
	}
	if item2.TotalQuantity != 5 {
		t.Errorf("Expected total quantity 5, got %d", item2.TotalQuantity)
	}
	if len(item2.Stores) != 1 {
		t.Errorf("Expected 1 store for item2, got %d", len(item2.Stores))
	}
}

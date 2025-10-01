package service_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/LudSkywalker/inventory-system/operator-service/app/dto"
	"github.com/LudSkywalker/inventory-system/operator-service/app/service"
	"github.com/LudSkywalker/inventory-system/operator-service/domain/entity"
)

type mockRepository struct {
	inventories map[string]*entity.GlobalInventory
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		inventories: make(map[string]*entity.GlobalInventory),
	}
}

func (m *mockRepository) Save(ctx context.Context, inventory *entity.GlobalInventory) error {
	key := inventory.StoreID + "-" + inventory.ItemID
	m.inventories[key] = inventory
	return nil
}

func (m *mockRepository) FindByItemAndStore(ctx context.Context, itemID, storeID string) (*entity.GlobalInventory, error) {
	key := storeID + "-" + itemID
	if inv, ok := m.inventories[key]; ok {
		return inv, nil
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockRepository) FindAll(ctx context.Context) ([]*entity.GlobalInventory, error) {
	var result []*entity.GlobalInventory
	for _, inv := range m.inventories {
		result = append(result, inv)
	}
	return result, nil
}

func (m *mockRepository) Delete(ctx context.Context, itemID, storeID string) error {
	key := storeID + "-" + itemID
	delete(m.inventories, key)
	return nil
}

func TestGlobalInventoryService(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepository()
	service := service.NewGlobalInventoryService(repo)

	// Test processing an inventory event
	t.Run("ProcessInventoryEvent", func(t *testing.T) {
		eventDTO := dto.GlobalInventoryDTO{
			ItemID:    "item1",
			StoreID:   "store1",
			Quantity:  10,
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		err := service.ProcessInventoryEvent(ctx, eventDTO)
		if err != nil {
			t.Errorf("ProcessInventoryEvent failed: %v", err)
		}

		// Verify the inventory was saved
		inventory, err := service.GetGlobalInventory(ctx, eventDTO.ItemID, eventDTO.StoreID)
		if err != nil {
			t.Errorf("GetGlobalInventory failed: %v", err)
		}

		if inventory.Quantity != eventDTO.Quantity {
			t.Errorf("Expected quantity %d, got %d", eventDTO.Quantity, inventory.Quantity)
		}
	})

	// Test getting all inventories
	t.Run("GetAllInventories", func(t *testing.T) {
		inventories, err := service.GetAllInventories(ctx)
		if err != nil {
			t.Errorf("GetAllInventories failed: %v", err)
		}

		if len(inventories) != 1 {
			t.Errorf("Expected 1 inventory, got %d", len(inventories))
		}
	})
}

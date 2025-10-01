package service

import (
	"context"

	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/infra/sqlite"
)

// InventoryService defines the interface for the inventory service.
type InventoryService interface {
	GetInventory(ctx context.Context, storeID, itemID string) (*sqlite.Inventory, error)
	UpdateInventory(ctx context.Context, inventory *sqlite.Inventory) error
	DeleteInventory(ctx context.Context, storeID, itemID string) error
	ListInventories(ctx context.Context) ([]*sqlite.Inventory, error)
}

// inventoryService implements InventoryService.
type inventoryService struct {
	repo *sqlite.SQLiteRepository
}

// NewInventoryService creates a new InventoryService.
func NewInventoryService(repo *sqlite.SQLiteRepository) InventoryService {
	return &inventoryService{
		repo: repo,
	}
}

func (s *inventoryService) GetInventory(ctx context.Context, storeID, itemID string) (*sqlite.Inventory, error) {
	return s.repo.GetInventory(ctx, storeID, itemID)
}

func (s *inventoryService) UpdateInventory(ctx context.Context, inventory *sqlite.Inventory) error {
	return s.repo.UpdateInventory(ctx, inventory)
}

func (s *inventoryService) DeleteInventory(ctx context.Context, storeID, itemID string) error {
	return s.repo.DeleteInventory(ctx, storeID, itemID)
}

func (s *inventoryService) ListInventories(ctx context.Context) ([]*sqlite.Inventory, error) {
	return s.repo.ListInventories(ctx)
}

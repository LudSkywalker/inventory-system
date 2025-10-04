package service

import (
	"context"
	"time"

	"github.com/LudSkywalker/inventory-system/internal/inventory/core/valueobject"
	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/app/dto"
	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/infra/sqlite"
)

// InventoryService defines the interface for the inventory service.
type InventoryService interface {
	GetInventory(ctx context.Context, storeID, itemID string) (*sqlite.Inventory, error)
	UpdateInventory(ctx context.Context, inventory *sqlite.Inventory) error
	DeleteInventory(ctx context.Context, storeID, itemID string) error
	ListInventories(ctx context.Context) ([]*sqlite.Inventory, error)
	UpdateStock(ctx context.Context, cmd dto.UpdateStockCommand) error
	GetStock(ctx context.Context, query dto.GetStockQuery) (*dto.InventoryDTO, error)
	DeleteStock(ctx context.Context, cmd dto.DeleteStockCommand) error
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

// Implement port.InventoryUseCase

func (s *inventoryService) UpdateStock(ctx context.Context, cmd dto.UpdateStockCommand) error {
	q, _ := valueobject.NewQuantity(cmd.Quantity)
	inventory := &sqlite.Inventory{
		ItemID:   cmd.ItemID,
		StoreID:  cmd.StoreID,
		Quantity: &q,
	}
	return s.repo.UpdateInventory(ctx, inventory)
}

func (s *inventoryService) GetStock(ctx context.Context, query dto.GetStockQuery) (*dto.InventoryDTO, error) {
	inv, err := s.repo.GetInventory(ctx, query.StoreID, query.ItemID)
	if err != nil {
		return nil, err
	}
	if inv == nil {
		return nil, nil
	}
	return &dto.InventoryDTO{
		ItemID:    inv.ItemID,
		StoreID:   inv.StoreID,
		Quantity:  inv.Quantity.Value(),
		UpdatedAt: inv.UpdatedAt.Format(time.RFC3339),
		Version:   inv.Version,
	}, nil
}

func (s *inventoryService) DeleteStock(ctx context.Context, cmd dto.DeleteStockCommand) error {
	return s.repo.DeleteInventory(ctx, cmd.StoreID, cmd.ItemID)
}

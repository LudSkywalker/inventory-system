package service

import (
	"context"
	"fmt"

	"github.com/LudSkywalker/inventory-system/internal/inventory/aggregator/domain"
	"github.com/LudSkywalker/inventory-system/internal/inventory/core/event"
)

type AggregatorService struct {
	repo domain.Repository
}

func NewAggregatorService(repo domain.Repository) *AggregatorService {
	return &AggregatorService{
		repo: repo,
	}
}

func (s *AggregatorService) ProcessInventoryEvent(ctx context.Context, event event.InventoryEvent) error {
	if event.Operation == event.OperationDelete {
		// For delete operations, we might want to keep the record but mark it as deleted
		// or implement a soft delete mechanism
		inventory := domain.NewGlobalInventory(event.ItemID, event.StoreID, 0)
		return s.repo.Save(ctx, inventory)
	}

	inventory, err := s.repo.FindByItemAndStore(ctx, event.ItemID, event.StoreID)
	if err != nil {
		// If not found, create new inventory
		inventory = domain.NewGlobalInventory(event.ItemID, event.StoreID, event.Quantity)
	} else {
		inventory.UpdateQuantity(event.Quantity)
	}

	if err := s.repo.Save(ctx, inventory); err != nil {
		return fmt.Errorf("saving global inventory: %w", err)
	}

	return nil
}

func (s *AggregatorService) GetAllInventories(ctx context.Context) ([]*domain.GlobalInventory, error) {
	return s.repo.FindAll(ctx)
}

func (s *AggregatorService) GetStoreInventory(ctx context.Context, storeID string) ([]*domain.GlobalInventory, error) {
	return s.repo.FindByStore(ctx, storeID)
}

func (s *AggregatorService) GetItemInventory(ctx context.Context, itemID string) ([]*domain.GlobalInventory, error) {
	return s.repo.FindByItem(ctx, itemID)
}

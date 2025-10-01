package app

import (
	"context"
	"time"

	"github.com/LudSkywalker/inventory-system/internal/inventory/core/entity"
	"github.com/LudSkywalker/inventory-system/internal/inventory/core/event"
	"github.com/LudSkywalker/inventory-system/internal/inventory/core/valueobject"
	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/domain"
)

// UpdateInventoryCommand represents a command to update inventory
type UpdateInventoryCommand struct {
	ItemID   string `json:"item_id"`
	StoreID  string `json:"store_id"`
	Quantity int    `json:"quantity"`
}

// InventoryService implements the application service for inventory management
type InventoryService struct {
	repo   domain.Repository
	events domain.EventPublisher
}

func NewInventoryService(repo domain.Repository, events domain.EventPublisher) *InventoryService {
	return &InventoryService{
		repo:   repo,
		events: events,
	}
}

func (s *InventoryService) UpdateInventory(ctx context.Context, cmd UpdateInventoryCommand) error {
	quantity, err := valueobject.NewQuantity(cmd.Quantity)
	if err != nil {
		return err
	}

	inventory, err := entity.NewInventory(cmd.ItemID, cmd.StoreID, quantity)
	if err != nil {
		return err
	}

	if err := s.repo.Save(ctx, inventory); err != nil {
		return err
	}

	// Publish event
	evt := &event.InventoryEvent{
		ItemID:    inventory.ItemID,
		StoreID:   inventory.StoreID,
		Quantity:  inventory.Quantity.Value(),
		Operation: event.OperationUpdate,
		Timestamp: inventory.UpdatedAt,
	}

	return s.events.PublishInventoryChange(ctx, evt)
}

func (s *InventoryService) GetInventory(ctx context.Context, itemID, storeID string) (*entity.Inventory, error) {
	return s.repo.FindByItemAndStore(ctx, itemID, storeID)
}

func (s *InventoryService) DeleteInventory(ctx context.Context, itemID, storeID string) error {
	if err := s.repo.Delete(ctx, itemID, storeID); err != nil {
		return err
	}

	evt := &event.InventoryEvent{
		ItemID:    itemID,
		StoreID:   storeID,
		Quantity:  0,
		Operation: event.OperationDelete,
		Timestamp: time.Now(),
	}

	return s.events.PublishInventoryChange(ctx, evt)
}

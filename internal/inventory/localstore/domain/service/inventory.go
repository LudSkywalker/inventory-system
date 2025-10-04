package service

import (
	"context"
	"fmt"
	"time"

	"github.com/LudSkywalker/inventory-system/internal/inventory/core/event"
	"github.com/LudSkywalker/inventory-system/internal/inventory/core/valueobject"
	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/domain"
)

type InventoryService struct {
	repo     domain.Repository
	eventBus domain.EventPublisher
}

func NewInventoryService(repo domain.Repository, eventBus domain.EventPublisher) *InventoryService {
	return &InventoryService{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (s *InventoryService) UpdateStock(ctx context.Context, itemID, storeID string, quantity int) error {
	qty, err := valueobject.NewQuantity(quantity)
	if err != nil {
		return fmt.Errorf("invalid quantity: %w", err)
	}

	inv, err := s.repo.Find(ctx, itemID, storeID)
	if err != nil {
		// Create new inventory if not found
		inv = domain.NewInventory(itemID, storeID, qty)
	} else {
		if err := inv.UpdateQuantity(qty); err != nil {
			return fmt.Errorf("updating quantity: %w", err)
		}
	}

	if err := s.repo.Save(ctx, inv); err != nil {
		return fmt.Errorf("saving inventory: %w", err)
	}

	// Publish event
	evt := &event.InventoryEvent{
		EventID:   "", // TODO: generate event ID
		ItemID:    inv.ItemID,
		StoreID:   inv.StoreID,
		Quantity:  inv.Quantity.Value(),
		Operation: event.OperationUpdate,
		Timestamp: time.Now(),
		Version:   inv.Version,
	}

	if err := s.eventBus.PublishInventoryChange(ctx, evt); err != nil {
		return fmt.Errorf("publishing event: %w", err)
	}

	return nil
}

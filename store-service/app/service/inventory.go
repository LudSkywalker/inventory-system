package service

import (
	"context"
	"fmt"

	"github.com/LudSkywalker/inventory-system/store-service/app/dto"
	output "github.com/LudSkywalker/inventory-system/store-service/app/port/output"
	"github.com/LudSkywalker/inventory-system/store-service/domain/entity"
	"github.com/LudSkywalker/inventory-system/store-service/domain/event"
	"github.com/LudSkywalker/inventory-system/store-service/domain/valueobject"
)

type inventoryService struct {
	repo     output.InventoryRepository
	eventBus output.EventPublisher
}

func NewInventoryService(repo output.InventoryRepository, eventBus output.EventPublisher) *inventoryService {
	return &inventoryService{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (s *inventoryService) UpdateStock(ctx context.Context, cmd dto.UpdateStockCommand) error {
	quantity, err := valueobject.NewQuantity(cmd.Quantity)
	if err != nil {
		return fmt.Errorf("invalid quantity: %w", err)
	}

	inventory, err := s.repo.Find(ctx, cmd.ItemID, cmd.StoreID)
	if err != nil {
		// If not found, create new inventory
		inventory, err = entity.NewInventory(cmd.ItemID, cmd.ItemName, cmd.StoreID, quantity)
		if err != nil {
			return fmt.Errorf("creating inventory: %w", err)
		}
	} else {
		inventory.ItemName = cmd.ItemName
		if err := inventory.UpdateQuantity(quantity); err != nil {
			return fmt.Errorf("updating quantity: %w", err)
		}
	}

	if err := s.repo.Save(ctx, inventory); err != nil {
		return fmt.Errorf("saving inventory: %w", err)
	}

	// Publish event
	evt := event.NewInventoryEvent(
		inventory.ItemID,
		inventory.ItemName,
		inventory.StoreID,
		inventory.Quantity.Value(),
		event.OperationUpdate,
	)

	if err := s.eventBus.PublishInventoryChange(ctx, evt); err != nil {
		return fmt.Errorf("publishing event: %w", err)
	}

	return nil
}

func (s *inventoryService) GetStock(ctx context.Context, query dto.GetStockQuery) (*dto.InventoryDTO, error) {
	inventory, err := s.repo.Find(ctx, query.ItemID, query.StoreID)
	if err != nil {
		return nil, fmt.Errorf("finding inventory: %w", err)
	}

	return &dto.InventoryDTO{
		ItemID:    inventory.ItemID,
		ItemName:  inventory.ItemName,
		StoreID:   inventory.StoreID,
		Quantity:  inventory.Quantity.Value(),
		UpdatedAt: inventory.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func (s *inventoryService) ListInventory(ctx context.Context) ([]*dto.InventoryDTO, error) {
	inventories, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing inventories: %w", err)
	}

	var result []*dto.InventoryDTO
	for _, inventory := range inventories {
		result = append(result, &dto.InventoryDTO{
			ItemID:    inventory.ItemID,
			ItemName:  inventory.ItemName,
			StoreID:   inventory.StoreID,
			Quantity:  inventory.Quantity.Value(),
			UpdatedAt: inventory.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return result, nil
}

func (s *inventoryService) SyncAllInventories(ctx context.Context) error {
	inventories, err := s.repo.List(ctx)
	if err != nil {
		return fmt.Errorf("listing inventories for sync: %w", err)
	}

	for _, inventory := range inventories {
		evt := event.NewInventoryEvent(
			inventory.ItemID,
			inventory.ItemName,
			inventory.StoreID,
			inventory.Quantity.Value(),
			event.OperationUpdate, // Use UPDATE for sync
		)

		if err := s.eventBus.PublishInventoryChange(ctx, evt); err != nil {
			return fmt.Errorf("publishing sync event for %s/%s: %w", inventory.ItemID, inventory.StoreID, err)
		}
	}

	return nil
}

func (s *inventoryService) DeleteStock(ctx context.Context, cmd dto.DeleteStockCommand) error {
	if err := s.repo.Delete(ctx, cmd.ItemID, cmd.StoreID); err != nil {
		return fmt.Errorf("deleting inventory: %w", err)
	}

	evt := event.NewInventoryEvent(
		cmd.ItemID,
		"", // ItemName not needed for delete
		cmd.StoreID,
		0,
		event.OperationDelete,
	)

	if err := s.eventBus.PublishInventoryChange(ctx, evt); err != nil {
		return fmt.Errorf("publishing event: %w", err)
	}

	return nil
}

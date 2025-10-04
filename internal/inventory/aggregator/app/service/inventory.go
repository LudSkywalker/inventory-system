package service

import (
	"context"
	"time"

	"github.com/LudSkywalker/inventory-system/internal/inventory/aggregator/app/dto"
	"github.com/LudSkywalker/inventory-system/internal/inventory/aggregator/domain"
	"github.com/LudSkywalker/inventory-system/internal/inventory/aggregator/domain/service"
	"github.com/LudSkywalker/inventory-system/internal/inventory/core/event"
)

type InventoryService struct {
	aggregator *service.AggregatorService
}

func NewInventoryService(aggregator *service.AggregatorService) *InventoryService {
	return &InventoryService{
		aggregator: aggregator,
	}
}

func (s *InventoryService) ProcessInventoryEvent(ctx context.Context, event event.InventoryEvent) error {
	return s.aggregator.ProcessInventoryEvent(ctx, event)
}

func (s *InventoryService) GetGlobalInventory(ctx context.Context, itemID, storeID string) (*dto.GlobalInventoryDTO, error) {
	inventories, err := s.aggregator.GetAllInventories(ctx)
	if err != nil {
		return nil, err
	}

	for _, inv := range inventories {
		if inv.ItemID == itemID && inv.StoreID == storeID {
			return &dto.GlobalInventoryDTO{
				ItemID:    inv.ItemID,
				StoreID:   inv.StoreID,
				Quantity:  inv.Quantity,
				UpdatedAt: inv.UpdatedAt.Format(time.RFC3339),
				Version:   inv.Version,
			}, nil
		}
	}

	return nil, nil
}

func (s *InventoryService) GetAllInventories(ctx context.Context) ([]*dto.GlobalInventoryDTO, error) {
	inventories, err := s.aggregator.GetAllInventories(ctx)
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(inventories), nil
}

func (s *InventoryService) GetStoreInventory(ctx context.Context, storeID string) ([]*dto.GlobalInventoryDTO, error) {
	inventories, err := s.aggregator.GetStoreInventory(ctx, storeID)
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(inventories), nil
}

func (s *InventoryService) GetItemInventory(ctx context.Context, itemID string) ([]*dto.GlobalInventoryDTO, error) {
	inventories, err := s.aggregator.GetItemInventory(ctx, itemID)
	if err != nil {
		return nil, err
	}

	return s.convertToDTO(inventories), nil
}

func (s *InventoryService) convertToDTO(inventories []*domain.GlobalInventory) []*dto.GlobalInventoryDTO {
	dtos := make([]*dto.GlobalInventoryDTO, len(inventories))
	for i, inv := range inventories {
		dtos[i] = &dto.GlobalInventoryDTO{
			ItemID:    inv.ItemID,
			StoreID:   inv.StoreID,
			Quantity:  inv.Quantity,
			UpdatedAt: inv.UpdatedAt.Format(time.RFC3339),
			Version:   inv.Version,
		}
	}
	return dtos
}

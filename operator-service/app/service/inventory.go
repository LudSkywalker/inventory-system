package service

import (
	"context"
	"fmt"
	"time"

	"github.com/LudSkywalker/inventory-system/operator-service/app/dto"
	"github.com/LudSkywalker/inventory-system/operator-service/app/port/output"
	"github.com/LudSkywalker/inventory-system/operator-service/domain/entity"
)

type globalInventoryService struct {
	repo output.GlobalInventoryRepository
}

func NewGlobalInventoryService(repo output.GlobalInventoryRepository) *globalInventoryService {
	return &globalInventoryService{
		repo: repo,
	}
}

func (s *globalInventoryService) ProcessInventoryEvent(ctx context.Context, eventDTO dto.GlobalInventoryDTO) error {
	inventory, err := s.repo.FindByItemAndStore(ctx, eventDTO.ItemID, eventDTO.StoreID)
	if err != nil {
		// If not found, create new inventory
		inventory = entity.NewGlobalInventory(eventDTO.ItemID, eventDTO.StoreID, eventDTO.Quantity)
	} else {
		inventory.UpdateQuantity(eventDTO.Quantity)
	}

	if err := s.repo.Save(ctx, inventory); err != nil {
		return fmt.Errorf("saving global inventory: %w", err)
	}

	return nil
}

func (s *globalInventoryService) GetGlobalInventory(ctx context.Context, itemID, storeID string) (*dto.GlobalInventoryDTO, error) {
	inventory, err := s.repo.FindByItemAndStore(ctx, itemID, storeID)
	if err != nil {
		return nil, fmt.Errorf("finding global inventory: %w", err)
	}

	return &dto.GlobalInventoryDTO{
		ItemID:    inventory.ItemID,
		StoreID:   inventory.StoreID,
		Quantity:  inventory.Quantity,
		UpdatedAt: inventory.UpdatedAt.Format(time.RFC3339),
		Version:   inventory.Version,
	}, nil
}

func (s *globalInventoryService) GetAllInventories(ctx context.Context) ([]*dto.GlobalInventoryDTO, error) {
	inventories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("finding all inventories: %w", err)
	}

	dtos := make([]*dto.GlobalInventoryDTO, len(inventories))
	for i, inventory := range inventories {
		dtos[i] = &dto.GlobalInventoryDTO{
			ItemID:    inventory.ItemID,
			StoreID:   inventory.StoreID,
			Quantity:  inventory.Quantity,
			UpdatedAt: inventory.UpdatedAt.Format(time.RFC3339),
			Version:   inventory.Version,
		}
	}

	return dtos, nil
}

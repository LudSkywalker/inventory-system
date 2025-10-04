package service

import (
	"context"
	"log"

	"github.com/LudSkywalker/inventory-system/backend-service/app/dto"
	"github.com/LudSkywalker/inventory-system/backend-service/infra/db"
)

type InventoryService struct {
	repo *db.HTTPRepository
}

func NewInventoryService(dbURL string) *InventoryService {
	return &InventoryService{
		repo: db.NewHTTPRepository(dbURL),
	}
}

func (s *InventoryService) GetGlobalInventory(ctx context.Context) ([]dto.InventoryDTO, error) {
	log.Println("GetGlobalInventory: Fetching all inventory")
	inventories, err := s.repo.FindAll(ctx)
	if err != nil {
		log.Printf("GetGlobalInventory: Error fetching inventory: %v", err)
		return nil, err
	}
	log.Printf("GetGlobalInventory: Found %d inventory items", len(inventories))

	var dtos []dto.InventoryDTO
	for _, inv := range inventories {
		dtos = append(dtos, *inv)
	}

	return dtos, nil
}

func (s *InventoryService) GetStoreInventory(ctx context.Context, storeID string) ([]dto.InventoryDTO, error) {
	inventories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var dtos []dto.InventoryDTO
	for _, inv := range inventories {
		if inv.StoreID == storeID {
			dtos = append(dtos, *inv)
		}
	}

	return dtos, nil
}

func (s *InventoryService) GetItemInventory(ctx context.Context, itemID string) ([]dto.InventoryDTO, error) {
	inventories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var dtos []dto.InventoryDTO
	for _, inv := range inventories {
		if inv.ItemID == itemID {
			dtos = append(dtos, *inv)
		}
	}

	return dtos, nil
}

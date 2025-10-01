package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LudSkywalker/inventory-system/internal/inventory/query/app/dto"
)

type InventoryQueryService struct {
	operatorURL string
}

func NewInventoryQueryService(operatorURL string) *InventoryQueryService {
	return &InventoryQueryService{
		operatorURL: operatorURL,
	}
}

func (s *InventoryQueryService) GetAllInventories(ctx context.Context) ([]*dto.InventoryDTO, error) {
	url := fmt.Sprintf("%s/api/v1/global-inventory", s.operatorURL)
	return s.doRequest(ctx, url)
}

func (s *InventoryQueryService) GetStoreInventory(ctx context.Context, storeID string) ([]*dto.InventoryDTO, error) {
	url := fmt.Sprintf("%s/api/v1/global-inventory/store/%s", s.operatorURL, storeID)
	return s.doRequest(ctx, url)
}

func (s *InventoryQueryService) GetItemInventory(ctx context.Context, itemID string) ([]*dto.InventoryDTO, error) {
	url := fmt.Sprintf("%s/api/v1/global-inventory/item/%s", s.operatorURL, itemID)
	return s.doRequest(ctx, url)
}

func (s *InventoryQueryService) GetInventory(ctx context.Context, query dto.InventoryQuery) (*dto.InventoryDTO, error) {
	if query.ItemID == "" || query.StoreID == "" {
		return nil, fmt.Errorf("both item ID and store ID are required")
	}

	url := fmt.Sprintf("%s/api/v1/global-inventory/%s/%s", s.operatorURL, query.StoreID, query.ItemID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var inventory dto.InventoryDTO
	if err := json.NewDecoder(resp.Body).Decode(&inventory); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &inventory, nil
}

func (s *InventoryQueryService) doRequest(ctx context.Context, url string) ([]*dto.InventoryDTO, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var inventories []*dto.InventoryDTO
	if err := json.NewDecoder(resp.Body).Decode(&inventories); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return inventories, nil
}

package db

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LudSkywalker/inventory-system/backend-service/app/dto"
)

type HTTPRepository struct {
	baseURL string
	client  *http.Client
}

func NewHTTPRepository(baseURL string) *HTTPRepository {
	return &HTTPRepository{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *HTTPRepository) FindAll(ctx context.Context) ([]*dto.InventoryDTO, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.baseURL+"/inventories", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var inventories []struct {
		ItemID    string    `json:"item_id"`
		StoreID   string    `json:"store_id"`
		Quantity  int       `json:"quantity"`
		UpdatedAt time.Time `json:"updated_at"`
		Version   int64     `json:"version"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&inventories); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	var dtos []*dto.InventoryDTO
	for _, inv := range inventories {
		dtos = append(dtos, &dto.InventoryDTO{
			ItemID:    inv.ItemID,
			StoreID:   inv.StoreID,
			Quantity:  inv.Quantity,
			UpdatedAt: inv.UpdatedAt.Format(time.RFC3339),
		})
	}

	return dtos, nil
}

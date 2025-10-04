package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	log.Printf("HTTPRepository: Making request to %s/inventories", r.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", r.baseURL+"/inventories", nil)
	if err != nil {
		log.Printf("HTTPRepository: Error creating request: %v", err)
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		log.Printf("HTTPRepository: Error sending request: %v", err)
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("HTTPRepository: Received response with status: %d", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Printf("HTTPRepository: Unexpected status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var inventories []struct {
		ItemID    string    `json:"item_id"`
		ItemName  string    `json:"item_name"`
		StoreID   string    `json:"store_id"`
		Quantity  int       `json:"quantity"`
		UpdatedAt time.Time `json:"updated_at"`
		Version   int64     `json:"version"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&inventories); err != nil {
		log.Printf("HTTPRepository: Error decoding response: %v", err)
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	log.Printf("HTTPRepository: Successfully decoded %d inventory items", len(inventories))
	var dtos []*dto.InventoryDTO
	for _, inv := range inventories {
		dtos = append(dtos, &dto.InventoryDTO{
			ItemID:    inv.ItemID,
			StoreID:   inv.StoreID,
			Quantity:  inv.Quantity,
			UpdatedAt: inv.UpdatedAt.Format(time.RFC3339),
		})
	}

	log.Printf("HTTPRepository: Returning %d inventory DTOs", len(dtos))
	return dtos, nil
}

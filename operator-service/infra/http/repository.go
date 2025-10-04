package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/LudSkywalker/inventory-system/operator-service/domain/entity"
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

func (r *HTTPRepository) Save(ctx context.Context, inventory *entity.GlobalInventory) error {
	reqBody := map[string]interface{}{
		"item_id":   inventory.ItemID,
		"item_name": inventory.ItemName,
		"store_id":  inventory.StoreID,
		"quantity":  inventory.Quantity,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.baseURL+"/inventories", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *HTTPRepository) FindByItemAndStore(ctx context.Context, itemID, storeID string) (*entity.GlobalInventory, error) {
	url := fmt.Sprintf("%s/inventories/%s/%s", r.baseURL, itemID, storeID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("inventory not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	var inventory entity.GlobalInventory
	if err := json.NewDecoder(resp.Body).Decode(&inventory); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Set IDs from URL parameters since response might not include them
	inventory.ItemID = itemID
	inventory.StoreID = storeID

	return &inventory, nil
}

func (r *HTTPRepository) FindAll(ctx context.Context) ([]*entity.GlobalInventory, error) {
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
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	var inventories []*entity.GlobalInventory
	if err := json.NewDecoder(resp.Body).Decode(&inventories); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return inventories, nil
}

func (r *HTTPRepository) Delete(ctx context.Context, itemID, storeID string) error {
	url := fmt.Sprintf("%s/inventories/%s/%s", r.baseURL, itemID, storeID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("inventory not found")
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

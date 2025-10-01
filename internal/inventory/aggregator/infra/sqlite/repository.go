package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LudSkywalker/inventory-system/internal/inventory/aggregator/domain"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Save(ctx context.Context, inventory *domain.GlobalInventory) error {
	query := `
		INSERT INTO global_inventories (item_id, store_id, quantity, updated_at, version)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(item_id, store_id) DO UPDATE SET
		quantity = ?,
		updated_at = ?,
		version = version + 1
	`

	_, err := r.db.ExecContext(ctx, query,
		inventory.ItemID,
		inventory.StoreID,
		inventory.Quantity,
		inventory.UpdatedAt,
		inventory.Version,
		inventory.Quantity,
		inventory.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("saving global inventory: %w", err)
	}

	return nil
}

func (r *SQLiteRepository) FindByItemAndStore(ctx context.Context, itemID, storeID string) (*domain.GlobalInventory, error) {
	query := `
		SELECT item_id, store_id, quantity, updated_at, version
		FROM global_inventories
		WHERE item_id = ? AND store_id = ?
	`

	inventory := &domain.GlobalInventory{}

	err := r.db.QueryRowContext(ctx, query, itemID, storeID).Scan(
		&inventory.ItemID,
		&inventory.StoreID,
		&inventory.Quantity,
		&inventory.UpdatedAt,
		&inventory.Version,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("inventory not found")
	}

	if err != nil {
		return nil, fmt.Errorf("querying global inventory: %w", err)
	}

	return inventory, nil
}

func (r *SQLiteRepository) FindAll(ctx context.Context) ([]*domain.GlobalInventory, error) {
	query := `
		SELECT item_id, store_id, quantity, updated_at, version
		FROM global_inventories
	`

	return r.queryInventories(ctx, query)
}

func (r *SQLiteRepository) FindByStore(ctx context.Context, storeID string) ([]*domain.GlobalInventory, error) {
	query := `
		SELECT item_id, store_id, quantity, updated_at, version
		FROM global_inventories
		WHERE store_id = ?
	`

	return r.queryInventories(ctx, query, storeID)
}

func (r *SQLiteRepository) FindByItem(ctx context.Context, itemID string) ([]*domain.GlobalInventory, error) {
	query := `
		SELECT item_id, store_id, quantity, updated_at, version
		FROM global_inventories
		WHERE item_id = ?
	`

	return r.queryInventories(ctx, query, itemID)
}

func (r *SQLiteRepository) queryInventories(ctx context.Context, query string, args ...interface{}) ([]*domain.GlobalInventory, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying inventories: %w", err)
	}
	defer rows.Close()

	var inventories []*domain.GlobalInventory
	for rows.Next() {
		inventory := &domain.GlobalInventory{}
		err := rows.Scan(
			&inventory.ItemID,
			&inventory.StoreID,
			&inventory.Quantity,
			&inventory.UpdatedAt,
			&inventory.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning inventory: %w", err)
		}
		inventories = append(inventories, inventory)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating inventories: %w", err)
	}

	return inventories, nil
}

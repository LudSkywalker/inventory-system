package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LudSkywalker/inventory-system/internal/inventory/core/valueobject"
)

type Inventory struct {
	StoreID   string
	ItemID    string
	Quantity  *valueobject.Quantity
	UpdatedAt time.Time
	Version   int64
}

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func InitDB(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS inventories (
			store_id TEXT NOT NULL,
			item_id TEXT NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 0,
			updated_at DATETIME NOT NULL,
			version INTEGER NOT NULL DEFAULT 1,
			PRIMARY KEY (store_id, item_id)
		);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating inventories table: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) GetInventory(ctx context.Context, storeID, itemID string) (*Inventory, error) {
	query := `
		SELECT store_id, item_id, quantity, updated_at, version
		FROM inventories
		WHERE store_id = ? AND item_id = ?
	`

	inv := &Inventory{}
	var quantity int
	err := r.db.QueryRowContext(ctx, query, storeID, itemID).Scan(
		&inv.StoreID,
		&inv.ItemID,
		&quantity,
		&inv.UpdatedAt,
		&inv.Version,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying inventory: %w", err)
	}

	q, _ := valueobject.NewQuantity(quantity)
	inv.Quantity = &q
	return inv, nil
}

func (r *SQLiteRepository) UpdateInventory(ctx context.Context, inventory *Inventory) error {
	query := `
		INSERT INTO inventories (store_id, item_id, quantity, updated_at, version)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(store_id, item_id) DO UPDATE SET
		quantity = excluded.quantity,
		updated_at = excluded.updated_at,
		version = version + 1
	`

	_, err := r.db.ExecContext(ctx, query,
		inventory.StoreID,
		inventory.ItemID,
		inventory.Quantity.Value(),
		inventory.UpdatedAt,
		inventory.Version,
	)

	if err != nil {
		return fmt.Errorf("error updating inventory: %w", err)
	}
	return nil
}

func (r *SQLiteRepository) DeleteInventory(ctx context.Context, storeID, itemID string) error {
	query := `DELETE FROM inventories WHERE store_id = ? AND item_id = ?`

	result, err := r.db.ExecContext(ctx, query, storeID, itemID)
	if err != nil {
		return fmt.Errorf("error deleting inventory: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("no inventory found for store %s and item %s", storeID, itemID)
	}

	return nil
}

func (r *SQLiteRepository) ListInventories(ctx context.Context) ([]*Inventory, error) {
	query := `
		SELECT store_id, item_id, quantity, updated_at, version
		FROM inventories
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying inventories: %w", err)
	}
	defer rows.Close()

	var inventories []*Inventory
	for rows.Next() {
		inv := &Inventory{}
		var quantity int
		if err := rows.Scan(&inv.StoreID, &inv.ItemID, &quantity, &inv.UpdatedAt, &inv.Version); err != nil {
			return nil, fmt.Errorf("error scanning inventory: %w", err)
		}
		q, _ := valueobject.NewQuantity(quantity)
		inv.Quantity = &q
		inventories = append(inventories, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return inventories, nil
}

func (r *SQLiteRepository) Delete(ctx context.Context, itemID, storeID string) error {
	query := "DELETE FROM inventories WHERE item_id = ? AND store_id = ?"

	result, err := r.db.ExecContext(ctx, query, itemID, storeID)
	if err != nil {
		return fmt.Errorf("deleting inventory: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("inventory not found")
	}

	return nil
}

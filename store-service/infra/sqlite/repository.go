package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LudSkywalker/inventory-system/store-service/domain/entity"
	"github.com/LudSkywalker/inventory-system/store-service/domain/valueobject"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Save(ctx context.Context, inventory *entity.Inventory) error {
	query := `
		INSERT INTO inventories (item_id, item_name, store_id, quantity, updated_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(item_id, store_id) DO UPDATE SET
		item_name = ?,
		quantity = ?,
		updated_at = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		inventory.ItemID,
		inventory.ItemName,
		inventory.StoreID,
		inventory.Quantity.Value(),
		inventory.UpdatedAt,
		inventory.ItemName,
		inventory.Quantity.Value(),
		inventory.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("saving inventory: %w", err)
	}

	return nil
}

func (r *SQLiteRepository) Find(ctx context.Context, itemID, storeID string) (*entity.Inventory, error) {
	query := `
		SELECT item_id, item_name, store_id, quantity, updated_at
		FROM inventories
		WHERE item_id = ? AND store_id = ?
	`

	var quantityValue int
	inventory := &entity.Inventory{}

	err := r.db.QueryRowContext(ctx, query, itemID, storeID).Scan(
		&inventory.ItemID,
		&inventory.ItemName,
		&inventory.StoreID,
		&quantityValue,
		&inventory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("inventory not found")
	}

	if err != nil {
		return nil, fmt.Errorf("querying inventory: %w", err)
	}

	quantity, err := valueobject.NewQuantity(quantityValue)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity in database: %w", err)
	}

	inventory.Quantity = quantity
	return inventory, nil
}

func (r *SQLiteRepository) List(ctx context.Context) ([]*entity.Inventory, error) {
	query := `
		SELECT item_id, item_name, store_id, quantity, updated_at
		FROM inventories
		ORDER BY updated_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying inventories: %w", err)
	}
	defer rows.Close()

	var inventories []*entity.Inventory
	for rows.Next() {
		var quantityValue int
		inventory := &entity.Inventory{}

		err := rows.Scan(
			&inventory.ItemID,
			&inventory.ItemName,
			&inventory.StoreID,
			&quantityValue,
			&inventory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning inventory: %w", err)
		}

		quantity, err := valueobject.NewQuantity(quantityValue)
		if err != nil {
			return nil, fmt.Errorf("invalid quantity in database: %w", err)
		}

		inventory.Quantity = quantity
		inventories = append(inventories, inventory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
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

// InitDB initializes the SQLite database with required tables
func InitDB(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS inventories (
			item_id TEXT NOT NULL,
			item_name TEXT NOT NULL,
			store_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			PRIMARY KEY (item_id, store_id)
		)
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("creating inventories table: %w", err)
	}

	return nil
}

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LudSkywalker/inventory-system/db-service/domain/entity"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Save(ctx context.Context, inventory *entity.GlobalInventory) error {
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

func (r *SQLiteRepository) FindByItemAndStore(ctx context.Context, itemID, storeID string) (*entity.GlobalInventory, error) {
	query := `
		SELECT item_id, store_id, quantity, updated_at, version
		FROM global_inventories
		WHERE item_id = ? AND store_id = ?
	`

	inventory := &entity.GlobalInventory{}
	var updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, itemID, storeID).Scan(
		&inventory.ItemID,
		&inventory.StoreID,
		&inventory.Quantity,
		&updatedAt,
		&inventory.Version,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("inventory not found")
	}

	if err != nil {
		return nil, fmt.Errorf("querying global inventory: %w", err)
	}

	inventory.UpdatedAt = updatedAt
	return inventory, nil
}

func (r *SQLiteRepository) FindAll(ctx context.Context) ([]*entity.GlobalInventory, error) {
	query := `
		SELECT item_id, store_id, quantity, updated_at, version
		FROM global_inventories
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying all global inventories: %w", err)
	}
	defer rows.Close()

	var inventories []*entity.GlobalInventory
	for rows.Next() {
		inventory := &entity.GlobalInventory{}
		var updatedAt time.Time

		err := rows.Scan(
			&inventory.ItemID,
			&inventory.StoreID,
			&inventory.Quantity,
			&updatedAt,
			&inventory.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning global inventory: %w", err)
		}

		inventory.UpdatedAt = updatedAt
		inventories = append(inventories, inventory)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating global inventories: %w", err)
	}

	return inventories, nil
}

func (r *SQLiteRepository) Delete(ctx context.Context, itemID, storeID string) error {
	query := "DELETE FROM global_inventories WHERE item_id = ? AND store_id = ?"

	result, err := r.db.ExecContext(ctx, query, itemID, storeID)
	if err != nil {
		return fmt.Errorf("deleting global inventory: %w", err)
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
		CREATE TABLE IF NOT EXISTS global_inventories (
			item_id TEXT NOT NULL,
			store_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			version INTEGER NOT NULL,
			PRIMARY KEY (item_id, store_id)
		)
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("creating global_inventories table: %w", err)
	}

	return nil
}
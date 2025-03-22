package database

import (
	"fmt"

	"github.com/ocmodi21/image-processing-service/internal/models"
)

func (db *PGClient) CreateStoreTable() error {
	query := `
		CREATE TABLE store (
			id VARCHAR(255) PRIMARY KEY,
			area_code VARCHAR(50),
			store_name VARCHAR(255)
		);
		CREATE INDEX store_id_idx ON store (id);
	`

	_, err := db.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to insert store: %w", err)
	}

	return nil
}

func (db *PGClient) InsertStore(store models.Store) error {
	query := `INSERT INTO store (id, area_code, store_name) VALUES ($1, $2, $3)`

	_, err := db.DB.Exec(query, store.ID, store.AreaCode, store.Name)
	if err != nil {
		return fmt.Errorf("failed to insert store: %w", err)
	}

	return nil
}

func (db *PGClient) GetStoreByID(storeID string) (*models.Store, error) {
	var store models.Store
	query := `SELECT id, area_code, store_name FROM store WHERE id = $1`

	err := db.DB.Get(&store, query, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	return &store, nil
}

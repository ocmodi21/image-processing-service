package storage

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/ocmodi21/image-processing-service/internal/database"
	"github.com/ocmodi21/image-processing-service/internal/models"
)

var (
	ErrStoreNotFound = errors.New("store not found")
)

// StoreStorage represents storage for store data
type StoreStorage struct {
	mu     sync.RWMutex
	stores map[string]*models.Store
	db     *database.PGClient
}

func NewStoreStorage(db *database.PGClient) *StoreStorage {
	return &StoreStorage{
		stores: make(map[string]*models.Store),
		db:     db,
	}
}

// LoadFromCSV loads store data from a CSV file
func (s *StoreStorage) LoadFromCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(record) >= 3 {
			store := &models.Store{
				AreaCode: record[0],
				Name:     record[1],
				ID:       record[2],
			}
			s.stores[store.ID] = store
		}
	}

	return nil
}

func (s *StoreStorage) GetStore(storeID string) (*models.Store, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	store, exists := s.stores[storeID]
	if !exists {
		return nil, ErrStoreNotFound
	}

	return store, nil
}

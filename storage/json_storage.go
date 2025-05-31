package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/imchukwu/finance-tracker/models"
)

type JSONStorage struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONStorage(filePath string) (*JSONStorage, error) {
	// Initialize with empty file if it doesn't exist
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(filePath, []byte("[]"), 0644); err != nil {
			return nil, err
		}
	}
	return &JSONStorage{filePath: filePath}, nil
}

func (j *JSONStorage) load() ([]*models.Transaction, error) {
	j.mu.RLock()
	defer j.mu.RUnlock()

	data, err := os.ReadFile(j.filePath)
	if err != nil {
		return nil, err
	}

	var transactions []*models.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (j *JSONStorage) save(transactions []*models.Transaction) error {
	j.mu.Lock()
	defer j.mu.Unlock()

	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(j.filePath, data, 0644)
}

func (j *JSONStorage) SaveTransaction(t *models.Transaction) error {
	transactions, err := j.LoadTransactions()
	if err != nil {
		return err
	}

	transactions = append(transactions, t)
	return j.save(transactions)
}

func (j *JSONStorage) LoadTransactions() ([]*models.Transaction, error) {
	return j.load()
}

func (j *JSONStorage) GetTransactionByID(id string) (*models.Transaction, error) {
	transactions, err := j.load()
	if err != nil {
		return nil, err
	}

	for _, t := range transactions {
		if t.ID == id {
			return t, nil
		}
	}

	return nil, errors.New("transaction not found")
}

func (j *JSONStorage) DeleteTransaction(id string) error {
	transactions, err := j.load()
	if err != nil {
		return err
	}

	filtered := make([]*models.Transaction, 0)
	for _, t := range transactions {
		if t.ID != id {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) == len(transactions) {
		return errors.New("transaction not found")
	}

	return j.save(filtered)
}

package storage

import "github.com/imchukwu/finance-tracker/models"

type Storage interface {
    SaveTransaction(t *models.Transaction) error
    LoadTransactions() ([]*models.Transaction, error)
	GetTransactionByID(id string) (*models.Transaction, error)
	DeleteTransaction(id string) error
}

// Temporary in-memory storage for now
type MemoryStorage struct {
	transactions []*models.Transaction
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		transactions: make([]*models.Transaction, 0),
	}
}

func (m *MemoryStorage) Save(t *models.Transaction) error {
	m.transactions = append(m.transactions, t)
	return nil
}

func (m *MemoryStorage) LoadAll() ([]*models.Transaction, error) {
	return m.transactions, nil
}


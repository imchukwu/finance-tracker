package models

import (
	"errors"
	"fmt"
	"time"
)

type Transaction struct {
	ID       string
	Date     time.Time
	Amount   float64
	Category string
	Notes    string
	Type     string
}

func NewTransaction(id string, date time.Time, amount float64,
	category, notes, transType string) *Transaction {
	return &Transaction{
		ID:       id,
		Date:     date,
		Amount:   amount,
		Category: category,
		Notes:    notes,
		Type:     transType,
	}
}

func (t *Transaction) IsExpense() bool {
	return t.Type == "expense"
}

func (t *Transaction) Display() string {
	return fmt.Sprintf("%s | %s | %s | %.2f | %s",
		t.ID,
		t.Date.Format("2006-01-02"),
		t.Type,
		t.Amount,
		t.Category)
}

func (t *Transaction) Validate() error {
    if t.ID == "" {
        return errors.New("transaction ID cannot be empty")
    }
    if t.Amount <= 0 {
        return errors.New("amount must be positive")
    }
    if t.Category == "" {
        return errors.New("category cannot be empty")
    }
    if t.Type != "income" && t.Type != "expense" {
        return errors.New("type must be either 'income' or 'expense'")
    }
    return nil
}
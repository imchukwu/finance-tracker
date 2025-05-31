package models

import (
	"fmt"
	"time"

	"github.com/imchukwu/finance-tracker/errors"
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
        return &errors.ValidationError{Field: "ID", Msg: "cannot be empty"}
    }
    if t.Amount <= 0 {
        return &errors.ValidationError{Field: "Amount", Msg: "must be positive"}
    }
    if t.Category == "" {
        return &errors.ValidationError{Field: "Category", Msg: "cannot be empty"}
    }
    if t.Type != "income" && t.Type != "expense" {
        return &errors.ValidationError{Field: "Type", Msg: "must be 'income' or 'expense'"}
    }
    return nil
}
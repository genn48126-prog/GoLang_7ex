package models

import "errors"

type Budget struct {
	Category    string  `json:"category"`
	LimitAmount float64 `json:"limit_amount"`
}

func (b *Budget) Validate() error {
	if b.Category == "" || b.LimitAmount <= 0 {
		return errors.New("некорректный бюджет")
	}
	return nil
}

type Transaction struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"` // YYYY-MM-DD
}

func (t *Transaction) Validate() error {
	if t.Amount == 0 || t.Category == "" || t.Date == "" {
		return errors.New("некорректная транзакция")
	}
	return nil
}
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"ledger/internal/models"
)

func AddTransaction(db *sql.DB, tx models.Transaction) (int, error) {
	if err := tx.Validate(); err != nil {
		return 0, err
	}

	var limit float64
	err := db.QueryRow("SELECT limit_amount FROM budgets WHERE category=$1", tx.Category).Scan(&limit)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == nil { // Ѕюджет существует
		var spent float64
		err = db.QueryRow("SELECT COALESCE(SUM(amount),0) FROM expenses WHERE category=$1", tx.Category).Scan(&spent)
		if err != nil {
			return 0, err
		}

		if spent+tx.Amount > limit {
			return 0, errors.New("budget exceeded")
		}
	}

	var id int
	query := `INSERT INTO expenses(amount, category, description, date) VALUES($1, $2, $3, $4) RETURNING id`
	err = db.QueryRow(query, tx.Amount, tx.Category, tx.Description, tx.Date).Scan(&id)
	return id, err
}

func ListTransactions(db *sql.DB) ([]models.Transaction, error) {
	rows, err := db.Query("SELECT id, amount, category, description, date FROM expenses ORDER BY date DESC, id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(&tx.ID, &tx.Amount, &tx.Category, &tx.Description, &tx.Date); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}
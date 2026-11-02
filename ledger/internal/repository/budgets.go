package repository

import (
	"database/sql"
	"ledger/internal/models"
)

func SetBudget(db *sql.DB, b models.Budget) error {
	if err := b.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO budgets(category, limit_amount) VALUES($1, $2)
		ON CONFLICT(category) DO UPDATE SET limit_amount = EXCLUDED.limit_amount`
	_, err := db.Exec(query, b.Category, b.LimitAmount)
	return err
}

func ListBudgets(db *sql.DB) ([]models.Budget, error) {
	rows, err := db.Query("SELECT category, limit_amount FROM budgets ORDER BY category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []models.Budget
	for rows.Next() {
		var b models.Budget
		if err := rows.Scan(&b.Category, &b.LimitAmount); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}
package handlers

import (
	"encoding/json"
	"net/http"

	"cashcraft/ledger/internal/cache"
	"cashcraft/ledger/internal/db"
	"cashcraft/ledger/internal/models"
	"cashcraft/ledger/internal/repository"
)

func SetBudgetHandler(w http.ResponseWriter, r *http.Request) {
	var b models.Budget
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := repository.SetBudget(db.DB, b); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(b)
}

func ListBudgetsHandler(w http.ResponseWriter, r *http.Request) {
	budgets, err := repository.ListBudgets(db.DB)
	if err != nil {
		http.Error(w, "Failed to fetch budgets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(budgets)
}

func AddTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var tx models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	id, err := repository.AddTransaction(db.DB, tx)
	if err != nil {
		// Ошибка превышения лимита
		if err.Error() == "budget exceeded" {
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}

func ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	txs, err := repository.ListTransactions(db.DB)
	if err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

func GetReportSummaryHandler(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		http.Error(w, "Parameters 'from' and 'to' are required", http.StatusBadRequest)
		return
	}

	summary, err := repository.GetReportSummary(r.Context(), db.DB, cache.Client, from, to)
	if err != nil {
		http.Error(w, "Failed to generate report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}
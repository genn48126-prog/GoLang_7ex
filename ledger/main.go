package main

import (
	"log"
	"net/http"

	"cashcraft/ledger/internal/cache"
	"cashcraft/ledger/internal/db"
	"cashcraft/ledger/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB()
	cache.InitRedis()

	r := mux.NewRouter()
	
	r.HandleFunc("/api/budgets", handlers.SetBudgetHandler).Methods("POST")
	r.HandleFunc("/api/budgets", handlers.ListBudgetsHandler).Methods("GET")
	
	r.HandleFunc("/api/transactions", handlers.AddTransactionHandler).Methods("POST")
	r.HandleFunc("/api/transactions", handlers.ListTransactionsHandler).Methods("GET")
	
	r.HandleFunc("/api/reports/summary", handlers.GetReportSummaryHandler).Methods("GET")

	// Ledger работает на порту 8081
	log.Println("Ledger service starting on :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
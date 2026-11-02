package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReportSummary struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

func GetReportSummary(ctx context.Context, db *sql.DB, rdb *redis.Client, from, to string) ([]ReportSummary, error) {
	cacheKey := "report:summary:" + from + ":" + to

	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var summary []ReportSummary
		if err := json.Unmarshal([]byte(cached), &summary); err == nil {
			log.Println("Cache HIT для отчета", cacheKey)
			return summary, nil
		}
	}

	log.Println("Cache MISS для отчета", cacheKey)

	query := `
		SELECT category, SUM(amount) as total 
		FROM expenses 
		WHERE date >= $1 AND date <= $2 
		GROUP BY category`
	rows, err := db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summary []ReportSummary
	for rows.Next() {
		var s ReportSummary
		if err := rows.Scan(&s.Category, &s.Total); err != nil {
			return nil, err
		}
		summary = append(summary, s)
	}

	jsonData, err := json.Marshal(summary)
	if err == nil {
		rdb.Set(ctx, cacheKey, jsonData, 30*time.Second)
	}

	return summary, nil
}
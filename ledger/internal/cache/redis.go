package cache

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := getEnv("REDIS_ADDR", "localhost:6379")
	dbStr := getEnv("REDIS_DB", "0")
	password := os.Getenv("REDIS_PASSWORD")

	db, _ := strconv.Atoi(dbStr)

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := Client.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	log.Println("Успешное подключение к Redis")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
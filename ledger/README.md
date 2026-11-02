# Ledger Service

Сервис учета финансов. Хранит данные в PostgreSQL, кэширует отчеты в Redis.

## Запуск инфраструктуры
```bash
docker run -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=cashapp -d postgres:15
docker run -p 6379:6379 --name cashcraft-redis -d redis:7-alpine
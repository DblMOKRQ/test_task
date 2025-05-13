# Trade Processing System

Микросервис для обработки торговых операций с HTTP API и фоновым воркером

## 📋 Требования

- Go 1.21+
- Docker 20.10+ и Docker Compose
- SQLite3 (для локальной разработки)

## 🚀 Быстрый старт с Docker

1. Соберите и запустите контейнеры:
```bash
docker-compose up --build
```
 2. Проверьте работоспособность:
 ```bash
 curl http://localhost:8080/healthz
```
## Локальная установка
1. Установите зависимости:
```bash
go mod download
```
2. Запустите сервер:
```
go run cmd/server/main.go -db data.db
```
3. В отдельном терминале запустите воркер:
```bash
go run cmd/worker/main.go -db data.db
```

## Примеры использования API
### Отправка сделки
```bash
curl -X POST http://localhost:8080/trades \
  -H "Content-Type: application/json" \
  -d '{
    "account": "test123",
    "symbol": "BTCUSD",
    "volume": 1.5,
    "open": 50000,
    "close": 51000,
    "side": "buy"
  }'
```
### Проверка здоровья
```bash
curl http://localhost:8080/healthz
```
### Получение статистики
```bash
curl http://localhost:8080/stats/test123
```
## Тестирование

Запуск всех тестов с проверкой race condition:
```bash
go test -race -covermode=atomic ./...
```
Проверка покрытия:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Docker окружение

### Структура сервисов

- `server`: HTTP API (порт 8080)
    
- `worker`: Фоновый обработчик сделок
    

### Переменные окружения

| Переменная | По умолчанию    | Описание                 |
| ---------- | --------------- | ------------------------ |
| DB_PATH    | /data/trades.db | Путь к файлу базы данных |

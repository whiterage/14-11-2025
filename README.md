## Web Links Status Service

Сервис на Go, реализующий требования тестового задания: принимает списки ссылок, асинхронно проверяет их доступность, присваивает задачам номера и по запросу формирует PDF‑отчёт по ранее отправленным наборам.

## Архитектура и ключевые решения
- **Go 1.25**, стандартная библиотека + `gofpdf` для работы с PDF.
- **In-memory репозиторий** с потокобезопасным счётчиком `links_num`.
- **Очередь заданий и пул воркеров**: `Service` кладёт задачи в канал, `WorkerPool` обрабатывает их и обновляет статусы.
- **HTTP API** на `net/http` + кастомные хендлеры (без сторонних фреймворков).
- **Graceful shutdown**: ловим `SIGINT/SIGTERM`, корректно останавливаем сервер и воркеров.

## Запуск
```bash
git clone git@github.com:whiterage/webserver_go.git
cd webserver_go
go run ./cmd/server
```
По умолчанию сервер слушает `http://localhost:8080`.

## Контрольный чек-лист (ручной прогон)
```bash
# 1. Юнит-тесты
go test ./...

# 2. Сборка бинаря
go build ./cmd/server

# 3. Запуск сервера (оставить в отдельном окне)
go run ./cmd/server

# 4. Health-check
curl -i http://localhost:8080/

# 5. Создание задачи со списком ссылок
curl -i -X POST http://localhost:8080/links \
  -H "Content-Type: application/json" \
  -d '{"links":["google.com","malformedlink.gg"]}'

# 6. Проверка статуса по links_num
curl -i http://localhost:8080/links/1

# 7. Получение PDF-отчёта
curl -X POST http://localhost:8080/links_list \
  -H "Content-Type: application/json" \
  -d '{"links_list":[1]}' \
  -o report.pdf
open report.pdf  # macOS

# 8. Graceful shutdown (в окне сервера)
Ctrl+C
```

## Примеры реальных ответов
```
> curl -i http://localhost:8080/
HTTP/1.1 200 OK
Content-Length: 2
Content-Type: text/plain; charset=utf-8
ok

> curl -i -X POST http://localhost:8080/links \
    -H "Content-Type: application/json" \
    -d '{"links":["google.com","malformedlink.gg"]}'
HTTP/1.1 201 Created
Content-Type: application/json
{"links":{"google.com":"pending","malformedlink.gg":"pending"},"links_num":1,"status":"pending"}

> curl -i http://localhost:8080/links/1
HTTP/1.1 200 OK
Content-Type: application/json
{"links":{"google.com":"available","malformedlink.gg":"not_available"},"links_num":1,"status":"done"}

> curl -X POST http://localhost:8080/links_list \
    -H "Content-Type: application/json" \
    -d '{"links_list":[1]}' \
    -o report.pdf
… файл report.pdf содержит таблицу со статусами ссылок и отметками времени проверки.
```

## API
### `POST /links`
```json
request:  { "links": ["google.com", "malformedlink.gg"] }
response: { "links": { "google.com": "pending", ... }, "links_num": 1, "status": "pending" }
```

### `GET /links/{links_num}`
Возвращает актуальные статусы по конкретному набору.

### `POST /links_list`
```json
request:  { "links_list": [1, 2] }
response: application/pdf (attachment)
```

## Технические детали
- **Пул воркеров**: размер задаётся в `cmd/server/main.go` (по умолчанию 4).
- **HTTPChecker** нормализует URL (добавляет `https://`, отбрасывает заведомо некорректные).
- **PDF отчёт**: включает заголовки, дату генерации, таблицы со ссылками, статусами и временем проверки.
- **Тесты**: базовые юнит‑тесты для вспомогательных функций (normalization, buildLinksMap). Команда запуска — `go test ./...`.

## Потенциальные улучшения
- Персистентное хранилище (SQLite/Postgres) для сохранения задач между перезапусками.
- Конфигурация через флаги/ENV (порт, таймауты, размер пула).
- Метрики и логирование в structured‑формате.


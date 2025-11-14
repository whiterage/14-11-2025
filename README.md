## Тестовое задание. Документация

Сервис на Go: принимает списки ссылок, асинхронно проверяет их доступность, присваивает задачам номера и по запросу формирует PDF‑отчёт по ранее отправленным наборам.

## Архитектура и ключевые решения
- **Go 1.25**, стандартная библиотека + `gofpdf` для работы с PDF.
- **In-memory репозиторий** с потокобезопасным счётчиком `links_num`.
- **Очередь заданий и пул воркеров**: `Service` кладёт задачи в канал, `WorkerPool` обрабатывает их и обновляет статусы.
- **HTTP API** на `net/http` + кастомные хендлеры (без сторонних фреймворков).
- **Graceful shutdown**: ловим `SIGINT/SIGTERM`, корректно останавливаем сервер и воркеров.

## Запуск
```bash
go run ./cmd/server
```
По умолчанию `http://localhost:8080`.

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
- **Персистентность**: задания и их статусы хранятся в `storage/tasks.json` (путь можно переопределить через `TASK_STORAGE_PATH`). При рестарте сервиса незавершённые задачи автоматически перезапускаются.
- **Graceful shutdown**: при `SIGINT/SIGTERM` сервер сначала завершает обработку HTTP‑запросов, затем ожидает, пока воркеры опустошат очередь задач; если лимит по времени превышен, воркеры принудительно отменяются.
- **Тесты**: помимо вспомогательных функций покрыта логика нормализации URL и работы с репозиторием. Команда запуска — `go test ./...`.
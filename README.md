# vk-forwarder

## Что делает

1. Поднимает HTTP-сервер на `SERVER_ADDR_VK_FORWARDER` и регистрирует `/webhook`/`/healthz`, используя маршрут из `internal/api`.
2. Проверяет подтверждение (`Confirmation`) и подпись (`Secret`) от VK, затем кидает сырые JSON-обновления в Kafka-топик `TOPIC_NAME_VK_UPDATES`.
3. Обрабатывает ошибки через zap-логгер и всегда отвечает `200 OK` для валидных webhook-запросов.
4. Логирует каждый `Send` с разделением по partition/offset, чтобы легко отслеживать публикацию событий.

## Запуск

1. Задайте конфигурацию в `.env` или экспортируйте переменные (`set -a && source .env && set +a`).
2. Соберите и запустите локально:
   ```bash
   go run ./cmd/vkforwarder
   ```
3. Или соберите Docker-образ и прокиньте переменные:
   ```bash
   docker build -t vk-forwarder .
   docker run --rm -e ... vk-forwarder
   ```

## Переменные окружения

Все переменные обязательны — сервис валидирует их и упадёт, если что-то отсутствует.

- `SERVER_ADDR_VK_FORWARDER` — адрес HTTP-сервера, по которому принимать webhook (например `0.0.0.0:8080`).
- `KAFKA_BOOTSTRAP_SERVERS_VALUE` — брокеры Kafka (`host:port[,host:port]`).
- `KAFKA_TOPIC_NAME_VK_UPDATES` — топик для сырых событий VK.
- `KAFKA_SASL_USERNAME` и `KAFKA_SASL_PASSWORD` — тайна SASL/Plain; для SCRAM их обязательно указывать.
- `VK_CONFIRMATION` — строка подтверждения, которую нужно отдать при регистрации webhook.
- `VK_SECRET` — ключ подписи `X-Hub-Signature`, используется для валидации `_sign`.
- `PATH_VK_FORWARDER_VK_WEB_HOOK` — путь, на котором слушает VK (например `/webhook`).
- `PATH_VK_FORWARDER_HEALTH_CHECK` — путь для `/healthz`, возвращает `200 OK`.

## Примечания

- Маршруты и проверка подписи реализованы в `internal/api`; любые изменения должны обновляться там.
- Kafka-продюсер фабричится через `internal/messaging`, использует `context.Background()` и `sarama.SyncProducer`.
- Если VK присылает `confirmation`, сервис отвечает строкой `VK_CONFIRMATION`, а не обрабатывает тело.
- Все логгирования идут через zap (`internal/logger`), ошибки приводят к `http.StatusBadRequest`.

# Insightly

AI-powered REST API для анализа CSV файлов. Пользователь загружает CSV, задаёт вопрос на естественном языке — сервис прогоняет данные через OpenAI и возвращает ответ, сохраняя историю запросов.

## Возможности

- Регистрация / авторизация по email и паролю, JWT access + refresh токены
- Загрузка CSV файлов, привязанных к пользователю
- AI-анализ содержимого файла по произвольному вопросу (OpenAI)
- История ранее заданных вопросов и ответов
- Rate limiting на защищённых роутах (Redis, 10 запросов/час на пользователя)
- Swagger-документация API

## Стек

- Go 1.26, [Gin](https://github.com/gin-gonic/gin)
- PostgreSQL (`sqlx`, миграции — `golang-migrate`)
- Redis (`go-redis/redis_rate` — rate limiting)
- OpenAI API (`openai-go`)
- JWT (`golang-jwt`)
- Swagger (`swaggo/swag`, `swaggo/gin-swagger`)
- Docker / Docker Compose

## Архитектура

Классическая слоистая архитектура:

```
handlers      → разбор HTTP-запроса, ответы
services      → бизнес-логика
repositories  → доступ к БД
models        → структуры данных
```

```
cmd/                    точка входа (main.go, роутинг)
internal/
  ai/                   клиент OpenAI
  db/                   подключение к Postgres
  errs/                 доменные ошибки
  handlers/             HTTP-хендлеры (users, files, queries, refresh_tokens)
  middlewares/          JWT-аутентификация, rate limiting, логирование
  models/               структуры User, Files, Queries, RefreshToken
  parser/                парсинг CSV в текст для AI
  repositories/         SQL-запросы
  services/             бизнес-логика
  mocks/                моки репозиториев/AI для юнит-тестов
migrations/             SQL-миграции (golang-migrate)
docs/                   сгенерированная swagger-документация
build/Dockerfile        сборка прод-образа
docker-compose.yml      app + postgres + redis
```

## Быстрый старт

### Требования

- Go 1.26+ (для локального запуска без Docker)
- Docker и Docker Compose (рекомендуемый способ)
- API-ключ OpenAI

### Настройка окружения

Скопируйте `.env.example` в `.env` и заполните значения:

```bash
cp .env.example .env
```

| Переменная      | Описание                                   |
|-----------------|---------------------------------------------|
| `OPENAI_API_KEY`| ключ OpenAI API                             |
| `JWT_SECRET`    | секрет для подписи access-токенов           |
| `POSTGRES_USER` | пользователь Postgres                       |
| `POSTGRES_DB`   | имя базы данных                             |
| `DB_PASSWORD`   | пароль Postgres                             |
| `POSTGRES_HOST` | хост Postgres (`db` внутри docker-сети)     |
| `POSTGRES_PORT` | порт Postgres (обычно `5432`)               |
| `POSTGRES_SSLMODE` | SSL-режим подключения (`disable` для локали) |
| `REDIS_URL`     | адрес Redis (`redis:6379` внутри docker-сети)|

### Запуск через Docker Compose

```bash
docker compose up -d --build
```

Поднимутся три контейнера: `app` (порт `8080`), `db` (Postgres, порт `5432`), `redis` (порт `6379`).

### Применение миграций

Миграции лежат в `migrations/` в формате [golang-migrate](https://github.com/golang-migrate/migrate). Применить их к поднятой базе:

```bash
migrate -path migrations -database "postgres://<user>:<password>@localhost:5432/<db>?sslmode=disable" up
```

### Локальный запуск без Docker

```bash
go run ./cmd
```

Postgres и Redis при этом должны быть доступны по адресам из `.env` (для локального запуска укажите `localhost` вместо `db`/`redis`).

## API-документация

После запуска сервиса Swagger UI доступен по адресу:

```
http://localhost:8080/swagger/index.html
```

Документация генерируется из аннотаций в коде хендлеров. После изменения аннотаций пересоздать `docs/`:

```bash
swag init -g cmd/main.go -o docs
```

### Основные эндпоинты

| Метод  | Путь               | Auth | Описание                              |
|--------|--------------------|:----:|----------------------------------------|
| POST   | `/auth/register`   |  —   | Регистрация нового пользователя        |
| POST   | `/auth/login`      |  —   | Аутентификация, выдача токенов         |
| POST   | `/auth/refresh`    |  —   | Обновление access-токена               |
| POST   | `/files/upload`    |  ✅  | Загрузка CSV файла                     |
| GET    | `/files`           |  ✅  | Список файлов текущего пользователя    |
| DELETE | `/files/{id}`      |  ✅  | Удаление файла                         |
| POST   | `/analyze`         |  ✅  | AI-анализ файла по вопросу             |
| GET    | `/analyze/history` |  ✅  | История запросов текущего пользователя |

Защищённые роуты требуют заголовок `Authorization: Bearer <access_token>`.

## Тесты

```bash
go test ./...
```

Юнит-тесты покрывают сервисный слой (`internal/services`) с использованием моков репозиториев и AI-клиента из `internal/mocks`. Файлы тестов в `internal/handlers` пока пустые заглушки.
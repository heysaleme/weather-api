# Weather API

REST API на Go для работы с погодой, пользовательскими городами и историей запросов.

## Что реализовано

- clean architecture с разделением на `handler`, `service`, `repository`
- dependency injection через конструкторы
- JWT authentication и роли `user` / `admin`
- middleware для auth, RBAC, request logging и recovery
- structured logging на `Uber Zap`
- unit tests для service и handler слоёв
- mock зависимости через `testify/mock`
- integration test для repository на SQLite in-memory
- негативные сценарии: `bad request`, `invalid JSON`, `not found`, `internal error`

## Архитектура

```text
Request
  -> LoggingMiddleware
  -> RecoveryMiddleware
  -> AuthMiddleware / RequireRole
  -> Handler
  -> Service
  -> Repository / External Client
```

Принципы:

- `handler` отвечает только за HTTP, DTO и status codes
- `service` содержит бизнес-логику
- `repository` отвечает за хранение данных
- зависимости передаются через интерфейсы и конструкторы
- внутренние ошибки не утекaют клиенту
- каждый запрос получает `request_id`, который попадает в response header и structured logs

## Structured logging

Все новые логи пишутся через `zap`.

Logging middleware логирует:

- `request_id`
- `method`
- `path`
- `status_code`
- `duration`

Recovery middleware:

- перехватывает panic
- пишет structured error log
- возвращает клиенту `500 internal server error`

Для внутренних ошибок API:

- клиент получает безопасное сообщение `internal server error`
- подробная ошибка остаётся только в логах

## Безопасность

- пароли хэшируются через `bcrypt`
- JWT secret хранится в environment variables
- токен подписывается `HS256`
- middleware валидирует подпись и срок жизни токена
- на защищённых маршрутах пользователь перепроверяется в repository
- `admin` не создаётся через публичную регистрацию
- `PasswordHash` не возвращается клиенту
- внутренние ошибки и panic не раскрываются наружу

## Конфигурация

Обязательная переменная:

```bash
export JWT_SECRET="super-secret-key"
```

Необязательные:

```bash
export HTTP_PORT="8080"
export BOOTSTRAP_ADMIN_EMAIL="admin@example.com"
export BOOTSTRAP_ADMIN_PASSWORD="very-strong-admin-password"
```

Если `BOOTSTRAP_ADMIN_*` заданы, при старте создаётся bootstrap admin.

## Запуск

```bash
go mod tidy
go run ./cmd
```

Сервер стартует на `http://localhost:8080`.

## Основные endpoints

### Auth

- `POST /auth/register`
- `POST /auth/login`

Пример регистрации:

```json
{
  "email": "user@example.com",
  "password": "strongpass123"
}
```

Пример логина:

```json
{
  "email": "user@example.com",
  "password": "strongpass123"
}
```

Response:

```json
{
  "access_token": "..."
}
```

### Protected routes

Заголовок:

```bash
Authorization: Bearer <token>
```

Пользовательские города:

- `POST /cities`
- `GET /cities`
- `DELETE /cities/{city_id}`

Погода пользователя:

- `GET /weather`
- `GET /weather/history`

Текущий пользователь:

- `GET /users/me`

### Admin routes

- `GET /users`
- `GET /users/{id}`
- `DELETE /users/{id}`

### Public weather routes

- `GET /weather/{city}`
- `GET /weather/country/{country}`
- `GET /weather/country/{country}/top`

## Примеры curl

Регистрация:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"strongpass123"}'
```

Логин:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"strongpass123"}'
```

Добавление города:

```bash
curl -X POST http://localhost:8080/cities \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"city":"Almaty"}'
```

Получение списка городов:

```bash
curl http://localhost:8080/cities \
  -H "Authorization: Bearer USER_TOKEN"
```

Погода и история:

```bash
curl http://localhost:8080/weather \
  -H "Authorization: Bearer USER_TOKEN"

curl http://localhost:8080/weather/history \
  -H "Authorization: Bearer USER_TOKEN"
```

## Тесты

В проекте есть:

- unit tests для `service`
- unit tests для `handler`
- unit tests для middleware
- integration test для SQLite repository

Основные команды:

```bash
GOCACHE=$(pwd)/.gocache go test ./...
GOCACHE=$(pwd)/.gocache go test -v ./...
GOCACHE=$(pwd)/.gocache go test -cover ./...
```

Если локальная среда не требует кастомный `GOCACHE`, можно запускать обычные команды:

```bash
go test ./...
go test -v ./...
go test -cover ./...
```

## Покрытие и проверяемые сценарии

Покрываются:

- happy path
- пустые данные
- invalid ID
- invalid JSON
- `not found`
- ошибки repository
- internal error
- panic recovery

Service layer покрыт выше требуемых `60%`.

## Структура проекта

```text
weather-api/
├── cmd/
│   ├── main.go
│   └── main_test.go
├── go.mod
├── go.sum
└── internal/
    ├── auth/
    ├── client/
    ├── config/
    ├── dto/
    ├── errs/
    ├── handler/
    ├── middleware/
    ├── model/
    ├── repository/
    └── service/
```

# Weather API

## Overview

REST API на Go для работы с погодой и пользовательскими данными.

Теперь сервис поддерживает:
- JWT аутентификацию
- роли `user` и `admin`
- защищённые маршруты
- хранение пользовательских городов
- историю погодных запросов
- разделение `model` и `dto`
- конфигурацию через `config`
- подготовленную структуру для unit tests service layer

Архитектура:

```text
Request -> AuthMiddleware -> Handler -> Service -> Repository
```

Принципы текущей структуры:
- `handler` работает только с HTTP и DTO
- `service` содержит бизнес-логику и зависит от interfaces
- `repository` отвечает только за доступ к данным
- `model` содержит доменные сущности
- `dto` содержит request/response контракты API
- `config` отвечает за загрузку env-конфигурации

## Run locally

Перед запуском задайте обязательные переменные окружения:

```bash
export JWT_SECRET="super-secret-key"
export BOOTSTRAP_ADMIN_EMAIL="admin@example.com"
export BOOTSTRAP_ADMIN_PASSWORD="very-strong-admin-password"
```

`BOOTSTRAP_ADMIN_*` необязательны. Если они заданы, при старте приложения создаётся начальный admin-аккаунт.

Запуск:

```bash
go mod tidy
go run cmd/main.go
```

Сервер стартует на `http://localhost:8080`.

## Security

- пароли хэшируются через `bcrypt`
- JWT secret хранится в env
- токен подписывается `HS256`
- в JWT есть `user_id`, `email`, `role`, `exp`
- middleware валидирует подпись и `exp`
- middleware перепроверяет пользователя в repository на каждом защищённом запросе
- admin роль не выдаётся через обычную регистрацию
- пароли не возвращаются в API и не хранятся в plain text

## Tests

В проекте уже есть unit tests для части критичных сценариев auth и middleware.

Запуск:

```bash
go test ./...
```

## Auth endpoints

### `POST /auth/register`

Регистрация пользователя.

Request:

```json
{
  "email": "user@example.com",
  "password": "strongpass123"
}
```

### `POST /auth/login`

Возвращает JWT access token.

Request:

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

## Protected endpoints

Для всех защищённых маршрутов нужен заголовок:

```bash
Authorization: Bearer <token>
```

### Cities

- `POST /cities`
- `GET /cities`
- `DELETE /cities/{city_id}`

Пример добавления города:

```json
{
  "city": "Almaty"
}
```

### Weather

- `GET /weather`
- `GET /weather/history`

`GET /weather` получает текущую погоду по всем городам текущего пользователя и сохраняет результат в историю.

### Current user

- `GET /users/me`

Возвращает данные текущего пользователя из JWT.

## Admin endpoints

Только для роли `admin`:

- `GET /users`
- `GET /users/{id}`
- `DELETE /users/{id}`

Admin создаётся только при старте приложения через `BOOTSTRAP_ADMIN_EMAIL` и `BOOTSTRAP_ADMIN_PASSWORD`, а не через публичную регистрацию.

## Expected status codes

- `401 Unauthorized` — нет токена, токен невалидный, истёк или пользователь больше не существует
- `403 Forbidden` — токен валидный, но роли недостаточно для admin endpoint
- `409 Conflict` — попытка зарегистрировать уже существующий email или добавить уже сохранённый город

## How to test

1. Зарегистрировать обычного пользователя:

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"strongpass123"}'
```

2. Залогиниться и получить `access_token`:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"strongpass123"}'
```

3. Добавить город:

```bash
curl -X POST http://localhost:8080/cities \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"city":"Almaty"}'
```

4. Получить города пользователя:

```bash
curl http://localhost:8080/cities \
  -H "Authorization: Bearer USER_TOKEN"
```

5. Получить текущую погоду и историю:

```bash
curl http://localhost:8080/weather \
  -H "Authorization: Bearer USER_TOKEN"

curl http://localhost:8080/weather/history \
  -H "Authorization: Bearer USER_TOKEN"
```

6. Проверить admin login и admin route:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"very-strong-admin-password"}'

curl http://localhost:8080/users \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

## Legacy public weather endpoints

Сохранены существующие публичные маршруты:

- `GET /weather/{city}`
- `GET /weather/country/{country}`
- `GET /weather/country/{country}/top`

## Project structure

```text
weather-api/
├── cmd/main.go
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

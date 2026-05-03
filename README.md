# Weather API

## Overview

REST API на Go для работы с погодой и пользовательскими данными.

Теперь сервис поддерживает:
- JWT аутентификацию
- роли `user` и `admin`
- защищённые маршруты
- хранение пользовательских городов
- историю погодных запросов

Архитектура:

```text
Request -> AuthMiddleware -> Handler -> Service -> Repository
```

## Run locally

Перед запуском задайте обязательные переменные окружения:

```bash
export JWT_SECRET="super-secret-key"
export ADMIN_EMAILS="admin@example.com"
```

`ADMIN_EMAILS` необязательна. Если email пользователя есть в этом списке, при регистрации он получает роль `admin`.

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
- пароли не возвращаются в API и не хранятся в plain text

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
    ├── errs/
    ├── handler/
    ├── middleware/
    ├── model/
    ├── repository/
    └── service/
```

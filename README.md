# Weather API

## Overview

Weather API is a REST service written in Go that provides current weather data by city and country.
The service integrates with external APIs to fetch real-time data and returns structured JSON responses.

The project follows a layered architecture:
- Handler — HTTP layer
- Service — business logic
- Client — external API interaction

---

## Getting Started

### Run locally

go mod tidy
go run cmd/main.go

Server will start on:

http://localhost:8080

---

## API Endpoints

### Get Weather by City

GET /weather/{city}

Description:
Returns current weather data for the specified city.

Flow:
1. Resolve city name to coordinates
2. Fetch weather data using coordinates
3. Add description and clothing recommendation

Example:
curl http://localhost:8080/weather/Almaty

Response:
{
  "city": "Almaty",
  "latitude": 43.25,
  "longitude": 76.95,
  "temperature": 25.4,
  "wind_speed": 3.2,
  "weather_code": 1,
  "time": "2026-04-16T12:00",
  "description": "Переменная облачность",
  "clothing": "Лёгкая одежда"
}

---

### Get Weather by Country

GET /weather/country/{country}

Description:
Returns weather data for a predefined list of cities within a country.

Behavior:
- Uses internal mapping of country → cities
- Fetches weather for each city
- Returns an array of results

Example:
curl http://localhost:8080/weather/country/Kazakhstan

---

### Get Top 3 Warmest Cities

GET /weather/country/{country}/top

Description:
Returns top 3 cities with the highest temperature.

Behavior:
- Retrieves all cities for the country
- Sorts them by temperature (descending)
- Returns up to 3 cities

Example:
curl http://localhost:8080/weather/country/Kazakhstan/top

---

## Clothing Recommendation

Temperature-based logic:

- < 5°C      → Тёплая одежда
- 5–15°C     → Куртка
- > 15°C     → Лёгкая одежда

---

## Architecture

Handler → Service → Client → External API

Handler (internal/handler):
- Parses HTTP requests
- Validates input
- Returns JSON responses

Service (internal/service):
- Contains business logic
- Aggregates and transforms data
- Adds derived fields (description, clothing)
- Handles sorting (top cities)

Client (internal/client):
- Performs HTTP requests to external APIs
- Parses responses

Model (internal/model):
- Defines response structures

---

## Project Structure

weather-api/
├── cmd/main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── client/
│   └── model/

---

## Error Handling

The API returns errors in JSON format:

{
  "error": "city not found"
}

---

## Technologies

- Go (net/http)
- Chi router
- Open-Meteo API
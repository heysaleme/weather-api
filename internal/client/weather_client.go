package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"weather-api/internal/errs"
)

type WeatherClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewWeatherClient(httpClient *http.Client) *WeatherClient {
	return &WeatherClient{
		httpClient: httpClient,
		baseURL:    "https://api.open-meteo.com/v1/forecast",
	}
}

type weatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
		Weathercode int     `json:"weathercode"`
		Time        string  `json:"time"`
	} `json:"current_weather"`
}

func (c *WeatherClient) GetCurrentWeather(ctx context.Context, lat, lon float64) (*weatherResponse, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%.4f", lat))
	q.Set("longitude", fmt.Sprintf("%.4f", lon))
	q.Set("current_weather", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.Upstream("weather request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.Upstream("weather status %d", resp.StatusCode)
	}

	var result weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errs.Upstream("decode weather response: %v", err)
	}

	return &result, nil
}

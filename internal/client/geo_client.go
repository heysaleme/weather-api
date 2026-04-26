package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"weather-api/internal/errs"
)

type GeoClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewGeoClient(httpClient *http.Client) *GeoClient {
	return &GeoClient{
		httpClient: httpClient,
		baseURL:    "https://geocoding-api.open-meteo.com/v1/search",
	}
}

type geoResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

func (c *GeoClient) GetCoordinates(ctx context.Context, city, countryCode string) (float64, float64, error) {
	city = strings.TrimSpace(city)
	if city == "" {
		return 0, 0, errs.InvalidInput("city is required")
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return 0, 0, err
	}

	q := u.Query()
	q.Set("name", city)
	q.Set("count", "1")
	if strings.TrimSpace(countryCode) != "" {
		q.Set("countryCode", strings.ToUpper(countryCode))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, 0, errs.Upstream("geocoding request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, errs.Upstream("geocoding status %d", resp.StatusCode)
	}

	var result geoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, errs.Upstream("decode geocoding response: %v", err)
	}

	if len(result.Results) == 0 {
		return 0, 0, errs.NotFound("city not found")
	}

	return result.Results[0].Latitude, result.Results[0].Longitude, nil
}

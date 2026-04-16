package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (c *GeoClient) GetCoordinates(ctx context.Context, city string) (float64, float64, error) {
	u, _ := url.Parse(c.baseURL)

	q := u.Query()
	q.Set("name", city)
	q.Set("count", "1")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var result geoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, err
	}

	if len(result.Results) == 0 {
		return 0, 0, fmt.Errorf("city not found")
	}

	return result.Results[0].Latitude, result.Results[0].Longitude, nil
}
